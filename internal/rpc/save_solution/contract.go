package save_solution

import (
	"context"
)

type Repository interface {
	SaveSolution(ctx context.Context, userID, competitionID int64) error
}
