package user_activity_total

import (
	"context"

	"github.com/varkis-ms/service-competition/internal/model"
)

type Repository interface {
	GetUserActivityTotal(ctx context.Context, userID int64) (model.UserActivityTotal, error)
}
