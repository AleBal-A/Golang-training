package main

import (
	"context"
	"flag"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"math"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
	pb "transmitter/transmitter"
)

type Stats struct {
	Count int     // Количество обработанных значений
	Mean  float64 // Среднее значение
	Std   float64 // Стандартное отклонение
	Sum   float64 // Сумма значений
	SumSq float64 // Сумма квадратов значений
}

func (s *Stats) AddValue(value float64) {
	s.Count++
	delta := value - s.Mean            // Разница между значением и текущим средним
	s.Mean += delta / float64(s.Count) // Обновление среднего значения
	s.Sum += value
	s.SumSq += value * value // Добавление квадрата значения к сумме квадратов
	if s.Count > 1 {         // Если количество значений больше одного, обновляется стандартное отклонение
		s.Std = math.Sqrt((s.SumSq / float64(s.Count)) - (s.Mean * s.Mean))
	}
}

type Anomaly struct {
	gorm.Model
	SessionID string `gorm:"index"`
	Frequency float64
	Mean      float64
	STD       float64
	Timestamp time.Time
}

func main() {
	k := flag.Float64("k", 3.0, "STD anomaly coefficient") // Коэффициент аномалии из командной строки
	flag.Parse()

	stats := &Stats{}

	valuesPool := sync.Pool{ // Пул для управления памятью
		New: func() interface{} {
			return new(float64)
		},
	}

	// Соединение с сервером
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	fmt.Printf("Successfully connect, k=%f \n", *k)
	defer conn.Close()

	client := pb.NewTransmitterServiceClient(conn)

	// Контекст с таймаутом в 600 секунд
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*600)
	defer cancel()

	// Создается запрос
	req := &pb.TransmitRequest{
		ClientId: "99",
	}

	ticker := time.NewTicker(10 * time.Second) // Таймер, срабатывающий каждые 10 секунд
	defer ticker.Stop()

	anomalyDetection := false // Флаг для переключения в режим обнаружения аномалий

	// Настройка соединения с базой данных PostgreSQL
	dsn := "host=localhost user=hakonoze password=h@konoze123 dbname=mydb port=5432 sslmode=disable TimeZone=Europe/Moscow"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Миграция схемы
	err = db.AutoMigrate(&Anomaly{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	go func() {
		for range ticker.C { // Каждые 10 секунд
			log.Printf("Processed values: %d, Predicted Mean: %f, Predicted STD: %f", stats.Count, stats.Mean, stats.Std)
		}
	}()

	// Поток данных от сервера
	stream, err := client.Transmit(ctx, req)
	if err != nil {
		log.Fatalf("could not transmit: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error receiving data: %v", err)
		}
		// Выводим полученные данные
		log.Printf("Recive frequency = %f", res.GetFrequency())
		value := valuesPool.Get().(*float64) // Значение из пула
		*value = res.GetFrequency()

		stats.AddValue(*value) // Добавляется значение в статистику

		if !anomalyDetection && stats.Count > 50 {
			anomalyDetection = true
			log.Println("Switching to Anomaly Detection stage")
		}

		if anomalyDetection { // Режим аномалий включен
			if math.Abs(*value-stats.Mean) > (*k * stats.Std) { // Является ли значение аномальным
				log.Printf("\n\n Anomaly detected: %f (Mean: %f, STD: %f \n\n)", *value, stats.Mean, stats.Std)

				// Сохранение аномалии в базу данных
				anomaly := &Anomaly{
					SessionID: res.SessionId,
					Frequency: *value,
					Mean:      stats.Mean,
					STD:       stats.Std,
					Timestamp: time.Time{},
				}
				if err := db.Create(&anomaly).Error; err != nil {
					log.Printf("failed to save anomaly: %v", err)
				}
			}
		}

		valuesPool.Put(value) // Возвращает значение в пул
	}
}
