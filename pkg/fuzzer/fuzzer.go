package fuzzer

import (
	"context"
)

// interface
type Fuzzer interface {
	Run(ctx context.Context)
}
