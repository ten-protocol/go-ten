package simulation

import "github.com/google/uuid"

// Utility function
func findDups(list []uuid.UUID) map[uuid.UUID]int {

	elementCount := make(map[uuid.UUID]int)

	for _, item := range list {
		// check if the item/element exist in the duplicate_frequency map

		_, exist := elementCount[item]

		if exist {
			elementCount[item] += 1 // increase counter by 1 if already in the map
		} else {
			elementCount[item] = 1 // else start counting from 1
		}
	}
	dups := make(map[uuid.UUID]int)
	for u, i := range elementCount {
		if i > 1 {
			dups[u] = i
		}
	}

	return dups
}
