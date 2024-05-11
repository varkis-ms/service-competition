package get_competition_info

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

func (h *Handler) Handle(ctx context.Context, in *pb.CompetitionInfoRequest, out *pb.CompetitionInfoResponse) error {
	log := h.log.With(
		slog.Int64("userID", in.UserID),
		slog.Int64("competitionID", in.CompetitionID),
	)

	comp, err := h.repo.GetCompetitionInfo(ctx, in.CompetitionID)
	if err != nil {
		log.Error("repo.GetCompetitionInfo failed", sl.Err(err))

		return err
	}

	maxScore, amountUsers, err := h.repo.GetAmountUsersAndMaxScore(ctx, in.CompetitionID)
	if err != nil {
		log.Error("repo.GetAmountUsersAndMaxScore failed", sl.Err(err))
	}

	out.CompetitionID = comp.CompetitionID
	out.Title = comp.Title
	out.Description = comp.Description
	out.DatasetTitle = comp.DatasetTitle
	out.DatasetDescription = comp.DatasetDescription
	out.AmountUsers = amountUsers
	out.MaximumScore = maxScore
	return nil
}
