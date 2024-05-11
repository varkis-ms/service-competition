package competition_list

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

func (h *Handler) Handle(ctx context.Context, in *pb.CompetitionListRequest, out *pb.CompetitionListResponse) error {
	log := h.log.With(
		slog.Int64("userID", in.UserID),
	)

	competitionList, err := h.repo.GetCompetitionList(ctx)
	if err != nil {
		log.Error("repo.GetCompetitionList failed", sl.Err(err))

		return err
	}

	for _, comp := range competitionList {
		competition := &pb.CompetitionInfo{
			CompetitionID: comp.CompetitionID,
			Title:         comp.Title,
			DatasetTitle:  comp.DatasetTitle,
		}

		out.CompetitionList = append(out.CompetitionList, competition)
	}

	return nil
}
