package main

import (
	"fmt"
)

func grayCode(n int) []int {
	size := 1 << n
	result := make([]int, size)
	for i := 0; i < size; i++ {
		result[i] = i ^ (i >> 1)
	}
	return result
}

func main() {
	var n int
	fmt.Print("Enter n: ")
	_, err := fmt.Scan(&n)
	if err != nil || n < 1 || n > 16 {
		fmt.Println("Invalid input. n must be between 1 and 16.")
		return
	}

	result := grayCode(n)
	fmt.Println(result)
}
