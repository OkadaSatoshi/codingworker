package main

import "testing"

func TestSum(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []int
		expected int
	}{
		{"empty slice", []int{}, 0},
		{"single element", []int{5}, 5},
		{"multiple elements", []int{1, 2, 3, 4, 5}, 15},
		{"negative numbers", []int{-1, -2, -3}, -6},
		{"mixed numbers", []int{-5, 10, -3, 8}, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sum(tt.numbers)
			if result != tt.expected {
				t.Errorf("Sum(%v) = %d, expected %d", tt.numbers, result, tt.expected)
			}
		})
	}
}

func TestAverage(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []int
		expected float64
	}{
		{"empty slice", []int{}, 0},
		{"single element", []int{10}, 10.0},
		{"multiple elements", []int{10, 20, 30}, 20.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Average(tt.numbers)
			if result != tt.expected {
				t.Errorf("Average(%v) = %.2f, expected %.2f", tt.numbers, result, tt.expected)
			}
		})
	}
}
