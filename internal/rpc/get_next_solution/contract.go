package get_next_solution

import (
	"context"
)

type Repository interface {
	GetNextSolutionInfo(ctx context.Context) (int64, int64, int64, error)
}
