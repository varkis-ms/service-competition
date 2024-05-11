package create_competition

import (
	"context"
	"log/slog"

	"github.com/varkis-ms/service-competition/internal/model"
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

func (h *Handler) Handle(ctx context.Context, in *pb.CompetitionCreateRequest, out *pb.CompetitionCreateResponse) error {
	log := h.log.With(
		slog.Int64("userID", in.UserID),
	)

	competitionID, err := h.repo.SaveCompetition(ctx, h.toModel(in))
	if err != nil {
		log.Error("repo.SaveCompetition failed", sl.Err(err))

		return err
	}

	out.CompetitionID = competitionID
	return nil
}

func (h *Handler) toModel(proto *pb.CompetitionCreateRequest) *model.Competition {
	return &model.Competition{
		UserID:             proto.UserID,
		Title:              proto.Title,
		Description:        proto.Description,
		DatasetTitle:       proto.DatasetTitle,
		DatasetDescription: proto.DatasetDescription,
	}
}
