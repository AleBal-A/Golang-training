package main

import "fmt"

type Present struct {
	Value int
	Size  int
}

// presents - Массив подарков, W - объём рюкзака
func grabPresents(presents []Present, W int) []Present {
	qty := len(presents)
	// Таблица для вычисления
	table := make([][]int, qty+1) // Количество подарков
	for i := range table {
		table[i] = make([]int, W+1) // Вместимость рюкзака
	}

	// Заполнение таблицы
	for i := 1; i <= qty; i++ {
		for w := 1; w <= W; w++ {
			if presents[i-1].Size <= w {
				// Если текущий подарок может поместиться в рюкзак
				table[i][w] = myMax(table[i-1][w], table[i-1][w-presents[i-1].Size]+presents[i-1].Value)
			} else {
				// Если текущий подарок не может поместиться в рюкзак
				table[i][w] = table[i-1][w]
			}
		}
	}
	// Вывод максимальной ценности, которую можно получить
	fmt.Println("Maximum value:", table[qty][W])

	// Теперь нам нужно восстановить, какие подарки были выбраны
	var selectedPresents []Present
	for i := qty; i > 0 && W > 0; i-- {
		if table[i][W] != table[i-1][W] {
			selectedPresents = append(selectedPresents, presents[i-1])
			W -= presents[i-1].Size
		}
	}

	return selectedPresents
}

func myMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {

	pp := []Present{
		{Value: 13, Size: 4},
		{Value: 5, Size: 1},
		{Value: 4, Size: 5},
		{Value: 3, Size: 1},
		{Value: 5, Size: 2},
	}

	fmt.Println(grabPresents(pp, 5))
}
