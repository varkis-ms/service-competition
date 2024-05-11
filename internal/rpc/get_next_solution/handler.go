package get_next_solution

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

func (h *Handler) Handle(ctx context.Context, out *pb.GetNextSolutionResponse) error {
	userID, compID, solutionID, err := h.repo.GetNextSolutionInfo(ctx)
	if err != nil {
		h.log.Error("repo.GetNextSolutionInfo failed", sl.Err(err))

		return err
	}

	out.CompetitionID = compID
	out.UserID = userID
	out.SolutionID = solutionID
	return nil
}
