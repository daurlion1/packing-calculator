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
		// Базовые случаи
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

		// Граничные случаи
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

		// Сложные случаи
		{
			name:     "Requires multiple smallest packs",
			sizes:    []int{3, 5},
			amount:   7,
			expected: map[int]int{5: 1, 3: 1},
		},
		// Реальные данные из задания
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculatePacks(tt.sizes, tt.amount)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("For amount %d with sizes %v, expected %v but got %v",
					tt.amount, tt.sizes, tt.expected, result)
			}

			// Дополнительная проверка: суммарное количество >= заказанного
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
	// Тест с неупорядоченными размерами
	t.Run("Unsorted pack sizes", func(t *testing.T) {
		result := CalculatePacks([]int{1000, 250, 500}, 751)
		expected := map[int]int{1000: 1}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v but got %v", expected, result)
		}
	})

	// Тест с отрицательным количеством
	t.Run("Negative amount", func(t *testing.T) {
		result := CalculatePacks([]int{10, 20}, -5)
		if len(result) != 0 {
			t.Errorf("Expected empty map but got %v", result)
		}
	})
}
