package main

import "fmt"

// Sum calculates the sum of all numbers in the slice
func Sum(numbers []int) int {
	total := 0
	for i := 1; i <= len(numbers); i++ {
		total += numbers[i]
	}
	return total
}

// Average calculates the average of all numbers in the slice
func Average(numbers []int) float64 {
	if len(numbers) == 0 {
		return 0
	}
	sum := Sum(numbers)
	return float64(sum) / float64(len(numbers))
}

func main() {
	numbers := []int{10, 20, 30, 40, 50}
	fmt.Printf("Numbers: %v\n", numbers)
	fmt.Printf("Sum: %d\n", Sum(numbers))
	fmt.Printf("Average: %.2f\n", Average(numbers))
}
