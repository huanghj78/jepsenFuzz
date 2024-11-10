package fuzzer

import (
	"context"
	"math/rand"
	"time"

	"github.com/huanghj78/jepsenFuzz/pkg/core"
)

type RandomFuzzer struct {
	nemesisGenerators core.NemesisGenerators
}

func (f *RandomFuzzer) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		time.Sleep(10)
		idx := rand.Intn(f.nemesisGenerators.GetGensLen())
		f.nemesisGenerators.SetIdx(idx)
	}
}

func NewRandomFuzzer(gens core.NemesisGenerators) *RandomFuzzer {
	return &RandomFuzzer{
		nemesisGenerators: gens,
	}
}
