package main

import (
	"container/heap"
	"errors"
	"fmt"
)

type Present struct {
	Value int
	Size  int
}

type PresentHeap []Present

func (h *PresentHeap) Len() int {
	return len(*h)
}

func (h *PresentHeap) Less(i, j int) bool {
	if (*h)[i].Value == (*h)[j].Value {
		return (*h)[i].Size < (*h)[j].Size
	}
	return (*h)[i].Value > (*h)[j].Value
}

func (h *PresentHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *PresentHeap) Push(x interface{}) {
	*h = append(*h, x.(Present))
}

func (h *PresentHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func getNCoolestPresents(presents []Present, n int) ([]Present, error) {
	if n <= 0 {
		return nil, errors.New("number of requested presents is equals to zero or negative")
	} else if n > len(presents) {
		return nil, errors.New("number of requested presents is more than available presents")
	}
	h := &PresentHeap{}
	heap.Init(h)

	for _, present := range presents {
		heap.Push(h, present)
	}

	coolestPresents := make([]Present, 0, n)
	for i := 0; i < n; i++ {
		coolestPresents = append(coolestPresents, heap.Pop(h).(Present))
	}

	return coolestPresents, nil
}

func main() {
	presents := []Present{
		{Value: 7, Size: 0},
		{Value: 4, Size: 5},
		{Value: 3, Size: 1},
		{Value: 5, Size: 2},
	}

	n := 0
	coolestPresents, err := getNCoolestPresents(presents, n)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Coolest presents:")
	for _, present := range coolestPresents {
		fmt.Println(present)
	}
}
