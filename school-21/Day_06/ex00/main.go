package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"os"
	"time"
)

func main() {
	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Определяем размеры изображения
	width := 300
	height := 300

	// Создаем новое изображение с указанными размерами
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Задаем цвет фона
	backgroundColor := color.RGBA{R: 35, G: 35, B: 30, A: 255} // Черный цвет
	draw.Draw(img, img.Bounds(), &image.Uniform{C: backgroundColor}, image.Point{}, draw.Src)

	// Рисуем фигуры или добавляем элементы логотипа
	// Пример: добавим белый круг в центр изображения
	circleColor := color.RGBA{R: uint8(rand.Intn(255)),
		G: uint8(rand.Intn(255)),
		B: uint8(rand.Intn(255)),
		A: 255} // Белый цвет

	radius := 100
	centerX, centerY := width/2, height/2
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				img.Set(centerX+x, centerY+y, circleColor)
			}
		}
	}

	// Сохраняем изображение в файл
	file, err := os.Create("amazing_logo.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}

	println("Логотип успешно создан: amazing_logo.png")
}
