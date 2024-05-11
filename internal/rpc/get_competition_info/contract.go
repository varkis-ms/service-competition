package get_competition_info

import (
	"context"

	"github.com/varkis-ms/service-competition/internal/model"
)

type Repository interface {
	GetCompetitionInfo(ctx context.Context, competitionID int64) (*model.Competition, error)
	GetAmountUsersAndMaxScore(ctx context.Context, competitionID int64) (float32, int64, error)
}
