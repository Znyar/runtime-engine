package main

import (
	"fmt"
	"math/rand"
	"time"
)

func generateRandomArray(size int) []int {
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = rand.Intn(1000000)
	}
	return arr
}

func bubbleSort(arr []int) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	const arraySize = 10000

	arr := generateRandomArray(arraySize)

	startTime := time.Now()

	bubbleSort(arr)

	elapsedTime := time.Since(startTime)

	fmt.Printf("Сортировка заняла: %v\n", elapsedTime)
}
