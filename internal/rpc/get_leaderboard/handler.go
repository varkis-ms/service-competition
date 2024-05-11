package get_leaderboard

import (
	"context"
	"log/slog"

	"github.com/varkis-ms/service-competition/internal/pkg/logger/sl"
	pb "github.com/varkis-ms/service-competition/internal/pkg/pb"
	"github.com/varkis-ms/service-competition/internal/storage"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=UserSaver

type Handler struct {
	repo Repository
	log  *slog.Logger
}

func New(
	repo storage.Repository,
	log *slog.Logger,
) *Handler {
	return &Handler{
		repo: repo,
		log:  log,
	}
}

func (h *Handler) Handle(ctx context.Context, in *pb.LeaderBoardRequest, out *pb.LeaderBoardResponse) error {
	log := h.log.With(slog.Int64(
		"competitionID", in.CompetitionID,
	))

	leaderboardList, err := h.repo.GetLeaderboard(ctx, in.CompetitionID)
	if err != nil {
		log.Error("repo.GetCompetitionInfo failed", sl.Err(err))

		return err
	}

	for _, leaderboard := range leaderboardList {
		leaderboardProto := &pb.LeaderBoard{
			UserID:  leaderboard.UserID,
			Score:   leaderboard.Score,
			AddedAt: leaderboard.AddedAt.String(),
		}

		out.LeaderBoardList = append(out.LeaderBoardList, leaderboardProto)
	}
	out.CompetitionID = in.CompetitionID
	return nil
}
