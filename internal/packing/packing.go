package packing

import (
	"sort"
)

// CalculatePacks determines the optimal combination of pack sizes to fulfill a customer's order.
// Rules applied:
// 1. Only whole packs can be used (no splitting).
// 2. Send the least number of *items* necessary to fulfill the order.
// 3. If multiple combinations satisfy rule 2, prefer the one with fewer *packs*.
//
// Params:
//   - sizes: slice of available pack sizes (e.g. [250, 500, 1000, ...])
//   - amount: number of items the customer ordered
//
// Returns:
//   - map[int]int where key = pack size, value = count of packs used
func CalculatePacks(sizes []int, amount int) map[int]int {
	// Ensure pack sizes are sorted from largest to smallest for optimization
	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))

	// Struct to store the best combination found so far
	best := struct {
		totalItems  int         // Total number of items (must be >= amount)
		numPacks    int         // Number of packs used
		combination map[int]int // Actual pack breakdown
	}{
		totalItems:  -1,
		numPacks:    -1,
		combination: nil,
	}

	// DFS function to explore all combinations recursively
	var dfs func(index, currentSum, packsUsed int, current map[int]int)
	dfs = func(index, currentSum, packsUsed int, current map[int]int) {
		// Base case: current combination satisfies the required amount
		if currentSum >= amount {
			// Check if this combination is better than the previous best
			if best.totalItems == -1 ||
				currentSum < best.totalItems ||
				(currentSum == best.totalItems && packsUsed < best.numPacks) {

				best.totalItems = currentSum
				best.numPacks = packsUsed

				// Deep copy the current combination to avoid mutation
				best.combination = make(map[int]int)
				for k, v := range current {
					best.combination[k] = v
				}
			}
			return
		}

		// Recursive case: try using more packs (can reuse same size)
		for i := index; i < len(sizes); i++ {
			size := sizes[i]
			current[size]++
			dfs(i, currentSum+size, packsUsed+1, current)
			current[size]--
			if current[size] == 0 {
				delete(current, size)
			}
		}
	}

	// Start DFS with empty combination
	dfs(0, 0, 0, make(map[int]int))

	// Return the best combination found, or empty if no solution
	if best.combination == nil {
		return map[int]int{}
	}
	return best.combination
}
