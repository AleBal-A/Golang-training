package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "transmitter/transmitter"
)

func main() {
	// Устанавливаем соединение с сервером
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTransmitterServiceClient(conn)

	// Устанавливаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Создаем запрос
	req := &pb.TransmitRequest{
		ClientId: "client123",
	}

	// Получаем поток данных от сервера
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
		log.Printf("Received data: session_id=%s, frequency=%f, timestamp=%s", res.GetSessionId(), res.GetFrequency(), res.GetTimestampUTC().AsTime().String())
	}
}
