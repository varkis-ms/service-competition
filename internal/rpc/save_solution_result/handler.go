package save_solution_result

import (
	"context"
	"log/slog"
	"time"

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

func (h *Handler) Handle(ctx context.Context, in *pb.SaveSolutionResultRequest) error {
	interval, err := time.ParseDuration(in.RunTime)
	if err != nil {
		h.log.Error("time.ParseDuration failed", sl.Err(err))

		return err
	}

	if err = h.repo.SaveSolutionInfo(ctx, in.SolutionID, in.Score, interval); err != nil {
		h.log.Error("repo.SaveSolutionInfo failed", sl.Err(err))

		return err
	}

	return nil
}
