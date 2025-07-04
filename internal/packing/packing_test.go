package packing

import (
	"reflect"
	"testing"
)

func TestCalculatePacks(t *testing.T) {
	tests := []struct {
		name     string
		sizes    []int
		amount   int
		expected map[int]int
	}{
		// Basic cases
		{
			name:     "Order exact pack size (single)",
			sizes:    []int{10, 20, 50},
			amount:   50,
			expected: map[int]int{50: 1},
		},
		{
			name:     "Order exact pack size (multiple)",
			sizes:    []int{10, 20, 50},
			amount:   70,
			expected: map[int]int{50: 1, 20: 1},
		},

		// Boundary cases
		{
			name:     "Order less than smallest pack",
			sizes:    []int{10, 20, 50},
			amount:   5,
			expected: map[int]int{10: 1},
		},
		{
			name:     "Order zero",
			sizes:    []int{10, 20, 50},
			amount:   0,
			expected: map[int]int{},
		},

		// Complex combinations
		{
			name:     "Requires multiple smallest packs",
			sizes:    []int{3, 5},
			amount:   7,
			expected: map[int]int{5: 1, 3: 1},
		},

		// Real examples from the task
		{
			name:     "Case from task: 251 items",
			sizes:    []int{250, 500, 1000, 2000, 5000},
			amount:   251,
			expected: map[int]int{500: 1},
		},
		{
			name:     "Case from task: 12001 items",
			sizes:    []int{250, 500, 1000, 2000, 5000},
			amount:   12001,
			expected: map[int]int{5000: 2, 2000: 1, 250: 1},
		},

		// High-volume tests
		{
			name:     "Large order: 100000 items",
			sizes:    []int{250, 500, 1000, 2000, 5000},
			amount:   100000,
			expected: map[int]int{5000: 20},
		},
		{
			name:     "Large uneven order: 100001 items",
			sizes:    []int{250, 500, 1000, 2000, 5000},
			amount:   100001,
			expected: map[int]int{5000: 20, 250: 1},
		},
		{
			name:     "Very large order: 999999 items",
			sizes:    []int{250, 500, 1000, 2000, 5000},
			amount:   999999,
			expected: map[int]int{5000: 200},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculatePacks(tt.sizes, tt.amount)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("For amount %d with sizes %v, expected %v but got %v",
					tt.amount, tt.sizes, tt.expected, result)
			}

			// Ensure total item count is >= requested amount
			total := 0
			for size, count := range result {
				total += size * count
			}
			if total < tt.amount {
				t.Errorf("Insufficient items: got %d for order %d", total, tt.amount)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("Unsorted pack sizes", func(t *testing.T) {
		result := CalculatePacks([]int{1000, 250, 500}, 751)
		expected := map[int]int{1000: 1}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v but got %v", expected, result)
		}
	})

	t.Run("Negative amount", func(t *testing.T) {
		result := CalculatePacks([]int{10, 20}, -5)
		if len(result) != 0 {
			t.Errorf("Expected empty map but got %v", result)
		}
	})
}
