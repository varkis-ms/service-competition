package save_solution_result

import (
	"context"
	"time"
)

type Repository interface {
	SaveSolutionInfo(ctx context.Context, solutionID int64, score float32, interval time.Duration) error
}
