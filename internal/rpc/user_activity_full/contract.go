package user_activity_full

import (
	"context"

	"github.com/varkis-ms/service-competition/internal/model"
)

type Repository interface {
	GetCompetitionInfoFull(ctx context.Context, userID int64) ([]model.CompetitionInfoFull, error)
	GetCompetitionInfoFullOwner(ctx context.Context, userID int64) ([]model.CompetitionInfoFullOwner, error)
}
