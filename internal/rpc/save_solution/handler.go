package save_solution

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

func (h *Handler) Handle(ctx context.Context, in *pb.SaveSolutionRequest) error {
	log := h.log.With(
		slog.Int64("userID", in.UserID),
	)

	if err := h.repo.SaveSolution(ctx, in.UserID, in.CompetitionID); err != nil {
		log.Error("repo.SaveSolution failed", sl.Err(err))

		return err
	}

	return nil
}
