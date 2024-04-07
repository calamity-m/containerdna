package heritage

import (
	"github.com/containers/image/v5/types"
	"github.com/sirupsen/logrus"
)

func ValidateChildParents(child []types.BlobInfo, parent ...[]types.BlobInfo) bool {
	parents := map[int][]types.BlobInfo{}
	for i, kv := range parent {
		if len(kv) > len(child) {
			// We cannot have a parent with a larger number of layers than childred
			logrus.Debug("One of the supplied parents has more layers than children")
			return false
		}
		parents[i] = kv
	}

	state := map[int]bool{}
	for i := range parent {
		state[i] = false
	}

	for i := range child {
		for j := range parents {
			if i >= len(parents[j]) {
				// Just skip if we are out of bounds for this individual parent, state is already captured
				continue
			}

			if child[i].Digest == parents[j][i].Digest {
				state[j] = true
			} else {
				state[j] = false
			}
		}
	}

	for _, v := range state {
		if !v {
			return false
		}
	}

	return true
}

/*
Lock step

	state := map[int]bool{}
	for i := range parent {
		state[i] = false
	}

	check := map[int]int{}
	for i := range parent {
		check[i] = 0
	}

	for i := range child {
		for j := range parents {

			if state[j] != true {
				// If our parent has not found a match yet, check the first index.

				if child[i].Digest == parents[j][0].Digest {
					// If check[j] is zero, we're still on the same layer concept

					// Flip our parent's state
					state[j] = true

					// Set where our parent's state was set
					check[j] = i

					// Say we have i = 80, for layer 80
					// When we go to layer 81 and find a new match we will have
					// 81 (i) - 80 (now check[j]) = 1
				}
			} else {
				// We found a match, and are currently matching. We need to check continued layers of parent

				if i-check[j] >= len(parents[j]) {
					// We are going to leave the bounds of our parent, so just stop. State will be what it is
					continue
				}

				if child[i].Digest == parents[j][i-check[j]].Digest {
					// Update our i
					check[j] = i
				} else {
					// Reset our state
					state[j] = false
					check[j] = 0
				}

			}
		}

	}


*/
