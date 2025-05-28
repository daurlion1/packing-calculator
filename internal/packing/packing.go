package packing

import (
	"sort"
)

func CalculatePacks(sizes []int, amount int) map[int]int {
	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))
	best := struct {
		totalItems  int
		numPacks    int
		combination map[int]int
	}{
		totalItems:  -1,
		numPacks:    -1,
		combination: nil,
	}

	var dfs func(index, currentSum, packsUsed int, current map[int]int)
	dfs = func(index, currentSum, packsUsed int, current map[int]int) {
		if currentSum >= amount {
			if best.totalItems == -1 ||
				currentSum < best.totalItems ||
				(currentSum == best.totalItems && packsUsed < best.numPacks) {

				best.totalItems = currentSum
				best.numPacks = packsUsed
				best.combination = make(map[int]int)
				for k, v := range current {
					best.combination[k] = v
				}
			}
			return
		}

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

	dfs(0, 0, 0, make(map[int]int))
	if best.combination == nil {
		return map[int]int{}
	}
	return best.combination
}
