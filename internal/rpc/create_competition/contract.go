package create_competition

import (
	"context"

	"github.com/varkis-ms/service-competition/internal/model"
)

type Repository interface {
	SaveCompetition(ctx context.Context, in *model.Competition) (int64, error)
}
