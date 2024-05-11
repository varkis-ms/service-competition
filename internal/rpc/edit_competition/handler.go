package edit_competition

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

func (h *Handler) Handle(ctx context.Context, in *pb.CompetitionEditRequest) error {
	log := h.log.With(
		slog.Int64("userID", in.UserID),
	)

	err := h.repo.EditCompetition(ctx, h.toModel(in))
	if err != nil {
		log.Error("repo.EditCompetition failed", sl.Err(err))

		return err
	}

	return nil
}

func (h *Handler) toModel(proto *pb.CompetitionEditRequest) *model.CompetitionEdit {
	return &model.CompetitionEdit{
		UserID:             proto.UserID,
		CompetitionID:      proto.CompetitionID,
		Title:              proto.Title,
		Description:        proto.Description,
		DatasetTitle:       proto.DatasetTitle,
		DatasetDescription: proto.DatasetDescription,
	}
}
