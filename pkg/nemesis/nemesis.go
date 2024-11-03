package nemesis

import (
	"math/rand"

	"github.com/huanghj78/jepsenFuzz/pkg/core"
)

func shuffleIndices(n int) []int {
	indices := make([]int, n)
	for i := 0; i < n; i++ {
		indices[i] = i
	}
	for i := len(indices) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		indices[i], indices[j] = indices[j], indices[i]
	}

	return indices
}

func init() {
	core.RegisterNemesis(kill{})
}
