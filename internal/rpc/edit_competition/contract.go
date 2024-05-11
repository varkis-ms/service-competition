package edit_competition

import (
	"context"

	"github.com/varkis-ms/service-competition/internal/model"
)

type Repository interface {
	EditCompetition(ctx context.Context, in *model.CompetitionEdit) error
}
