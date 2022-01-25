package simulation

import "github.com/google/uuid"

// Utility function
func findDups(list []uuid.UUID) map[uuid.UUID]int {

	duplicate_frequency := make(map[uuid.UUID]int)

	for _, item := range list {
		// check if the item/element exist in the duplicate_frequency map

		_, exist := duplicate_frequency[item]

		if exist {
			duplicate_frequency[item] += 1 // increase counter by 1 if already in the map
		} else {
			duplicate_frequency[item] = 1 // else start counting from 1
		}
	}
	dups := make(map[uuid.UUID]int)
	for u, i := range duplicate_frequency {
		if i > 1 {
			dups[u] = i
		}
	}

	return dups
}
