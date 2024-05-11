package storage

import (
	"context"
	"time"

	"github.com/varkis-ms/service-competition/internal/model"
)

// Repository описывает операции на уровне хранилища
type Repository interface {
	SaveCompetition(ctx context.Context, in *model.Competition) (int64, error)
	EditCompetition(ctx context.Context, in *model.CompetitionEdit) error
	GetCompetitionList(ctx context.Context) (model.CompetitionList, error)
	GetCompetitionInfo(ctx context.Context, competitionID int64) (*model.Competition, error)
	GetLeaderboard(ctx context.Context, competitionID int64) (model.LeaderBoardList, error)
	GetAmountUsersAndMaxScore(ctx context.Context, competitionID int64) (float32, int64, error)
	GetUserActivityTotal(ctx context.Context, userID int64) (model.UserActivityTotal, error)
	GetCompetitionInfoFull(ctx context.Context, userID int64) ([]model.CompetitionInfoFull, error)
	GetCompetitionInfoFullOwner(ctx context.Context, userID int64) ([]model.CompetitionInfoFullOwner, error)
	SaveSolutionInfo(ctx context.Context, solutionID int64, score float32, interval time.Duration) error
	GetNextSolutionInfo(ctx context.Context) (int64, int64, int64, error)
	//GetCurrentSolve(ctx context.Context) (string, string, int, error)
	//SaveSolveDataToDb(ctx context.Context, userId, competitionID string) error
	//SaveScoreToDb(ctx context.Context, queueId, score int) error
}
