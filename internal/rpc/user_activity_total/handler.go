package user_activity_total

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

func (h *Handler) Handle(ctx context.Context, in *pb.UserActivityTotalRequest, out *pb.UserActivityTotalResponse) error {
	log := h.log.With(slog.Int64(
		"userID", in.UserID,
	))

	userActivityTotal, err := h.repo.GetUserActivityTotal(ctx, in.UserID)
	if err != nil {
		log.Error("repo.GetUserActivityTotal failed", sl.Err(err))

		return err
	}

	out.TotalTime = userActivityTotal.TotalTime
	out.TotalAttempts = userActivityTotal.TotalAttempts
	out.TotalCompetitions = userActivityTotal.TotalCompetitions
	out.TotalOwnerCompetitions = userActivityTotal.TotalOwnerCompetitions
	return nil
}
