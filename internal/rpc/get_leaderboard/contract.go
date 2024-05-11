package get_leaderboard

import (
	"context"

	"github.com/varkis-ms/service-competition/internal/model"
)

type Repository interface {
	GetLeaderboard(ctx context.Context, competitionID int64) (model.LeaderBoardList, error)
}
