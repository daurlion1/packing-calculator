package packing

import (
	"sort"
)

// CalculatePacks determines the optimal combination of pack sizes to fulfill a customer's order.
// The algorithm follows three strict rules:
//  1. Only whole packs may be used (no splitting).
//  2. Send the least number of items possible, but not fewer than the ordered amount.
//  3. If multiple combinations send the same number of items, choose the one with fewer packs.
//
// Parameters:
//   - sizes: available pack sizes (e.g., [250, 500, 1000, 2000, 5000])
//   - amount: number of items ordered by the customer
//
// Returns:
//   - A map[int]int where key = pack size, and value = number of packs used
func CalculatePacks(sizes []int, amount int) map[int]int {
	// Sort pack sizes in ascending order for easier traversal in DP
	sort.Ints(sizes)

	type dpState struct {
		index  int // Current index in sizes slice
		amount int // Remaining amount to fulfill
	}

	type dpResult struct {
		totalItems int         // Total items used (must be >= order amount)
		numPacks   int         // Number of packs used
		comb       map[int]int // Pack size -> count mapping
	}

	// Memoization cache for DP results
	cache := make(map[dpState]*dpResult)

	// Recursive DP function to explore combinations from the current index
	var findOptimal func(index, remaining int) *dpResult
	findOptimal = func(index, remaining int) *dpResult {
		state := dpState{index, remaining}
		if val, ok := cache[state]; ok {
			return val
		}

		// If we've considered all sizes and still have items to fulfill
		if index == len(sizes) {
			return &dpResult{totalItems: -1}
		}

		best := &dpResult{totalItems: -1}
		currentPack := sizes[index]

		// Max packs of current size needed to (over)fulfill the remaining amount
		maxCount := (remaining + currentPack - 1) / currentPack

		// Try using 0 to maxCount packs of current size
		for count := 0; count <= maxCount; count++ {
			remainingAfter := remaining - count*currentPack
			if remainingAfter < 0 {
				remainingAfter = 0
			}

			subResult := findOptimal(index+1, remainingAfter)

			// If subResult is invalid and we still have a gap to fill — skip
			if subResult.totalItems == -1 && remainingAfter > 0 {
				continue
			}

			// Compute total items and pack count for current combination
			totalItems := count * currentPack
			numPacks := count

			if subResult.totalItems != -1 {
				totalItems += subResult.totalItems
				numPacks += subResult.numPacks
			}

			// Check if this combination is better than the best so far
			if best.totalItems == -1 ||
				totalItems < best.totalItems ||
				(totalItems == best.totalItems && numPacks < best.numPacks) {

				// Deep copy the combination to avoid mutation
				best.totalItems = totalItems
				best.numPacks = numPacks
				best.comb = make(map[int]int)

				if subResult.totalItems != -1 {
					for k, v := range subResult.comb {
						best.comb[k] = v
					}
				}
				if count > 0 {
					best.comb[currentPack] += count
				}
			}
		}

		cache[state] = best
		return best
	}

	// Start recursion from the first pack size and the full order amount
	result := findOptimal(0, amount)
	if result.totalItems == -1 {
		// No valid combination found
		return map[int]int{}
	}
	return result.comb
}
