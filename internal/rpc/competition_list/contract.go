package competition_list

import (
	"context"

	"github.com/varkis-ms/service-competition/internal/model"
)

type Repository interface {
	GetCompetitionList(context.Context) (model.CompetitionList, error)
}
