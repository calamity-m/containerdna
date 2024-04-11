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
