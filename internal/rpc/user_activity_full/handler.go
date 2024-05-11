package user_activity_full

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
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

func (h *Handler) Handle(ctx context.Context, in *pb.UserActivityFullRequest, out *pb.UserActivityFullResponse) error {
	log := h.log.With(slog.Int64(
		"userID", in.UserID,
	))

	competitionInfoFull, err := h.repo.GetCompetitionInfoFull(ctx, in.UserID)
	if err != nil {
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				log.Error("repo.GetCompetitionInfoFull failed", sl.Err(err))
				return err
			}
		}
	}

	competitionInfoFullOwner, err := h.repo.GetCompetitionInfoFullOwner(ctx, in.UserID)
	if err != nil {
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				log.Error("repo.GetCompetitionInfoFullOwner failed", sl.Err(err))
				return err
			}
		}
	}

	h.toProto(out, competitionInfoFull, competitionInfoFullOwner)
	return nil
}

func (h *Handler) toProto(
	proto *pb.UserActivityFullResponse,
	modelFull []model.CompetitionInfoFull,
	modelFullOwner []model.CompetitionInfoFullOwner,
) {
	var member []*pb.CompetitionInfoFull
	for _, itemFull := range modelFull {
		member = append(member, &pb.CompetitionInfoFull{
			CompetitionID: itemFull.CompetitionID,
			Title:         itemFull.Title,
			DatasetTitle:  itemFull.DatasetTitle,
			Score:         itemFull.Score,
			AddedAt:       itemFull.AddedAt.String(),
		})
	}

	var owner []*pb.CompetitionInfoFullOwner
	for _, itemFullOwner := range modelFullOwner {
		owner = append(owner, &pb.CompetitionInfoFullOwner{
			CompetitionID: itemFullOwner.CompetitionID,
			Title:         itemFullOwner.Title,
			DatasetTitle:  itemFullOwner.DatasetTitle,
			AmountUsers:   itemFullOwner.AmountUsers,
			AddedAt:       itemFullOwner.AddedAt.String(),
		})
	}

	proto.Member = member
	proto.Owner = owner
}
