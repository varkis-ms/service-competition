package storage

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/varkis-ms/service-competition/internal/model"
)

func (r *Storage) GetUserActivityTotal(ctx context.Context, userID int64) (model.UserActivityTotal, error) {
	sql, args, _ := r.Builder.
		Select("SUM(run_time)",
			"COALESCE(COUNT(DISTINCT competition_id), 0)",
			"COALESCE(COUNT(*), 0)",
			"COALESCE(COUNT(DISTINCT lb.competition_id), 0)").
		From("leaderboard lb").
		Join("competition comp ON lb.competition_id = comp.id").
		Where(sq.Eq{"lb.user_id": userID}).
		ToSql()

	var res model.UserActivityTotal
	var totalRunTime *time.Duration
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&totalRunTime,
		&res.TotalAttempts,
		&res.TotalCompetitions,
		&res.TotalOwnerCompetitions,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.UserActivityTotal{}, nil
		}

		return model.UserActivityTotal{}, err
	}

	if totalRunTime != nil {
		res.TotalTime = totalRunTime.String()
	}
	return res, nil
}

func (r *Storage) GetCompetitionInfoFull(ctx context.Context, userID int64) ([]model.CompetitionInfoFull, error) {
	sql, args, _ := r.Builder.
		Select("COALESCE(lb.competition_id, 0)",
			"COALESCE(comp.title, '')",
			"COALESCE(comp.dataset_title, '')",
			"COALESCE(lb.score, 0)",
			"CASE WHEN comp.created_at IS NOT NULL THEN comp.created_at ELSE '0001-01-01 00:00:00' END").
		From("leaderboard lb").
		Join("competition comp ON lb.competition_id = comp.id").
		Where(sq.Eq{"lb.user_id": userID}).
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}
	defer rows.Close()

	var res []model.CompetitionInfoFull
	for rows.Next() {
		var item model.CompetitionInfoFull
		err := rows.Scan(&item.CompetitionID, &item.Title, &item.DatasetTitle, &item.Score, &item.AddedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Storage) GetCompetitionInfoFullOwner(ctx context.Context, userID int64) ([]model.CompetitionInfoFullOwner, error) {
	sql, args, _ := r.Builder.
		Select("comp.id AS competitionID",
			"comp.title",
			"comp.dataset_title",
			"COUNT(lb.user_id) AS amountUsers",
			"comp.created_at AS created_at").
		From("competition comp").
		Join("leaderboard lb ON comp.id = lb.competition_id").
		Where("comp.user_id = ?", userID).
		GroupBy("comp.id, comp.title, comp.dataset_title, comp.created_at").
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}
	defer rows.Close()

	var res []model.CompetitionInfoFullOwner
	for rows.Next() {
		var item model.CompetitionInfoFullOwner
		err := rows.Scan(&item.CompetitionID, &item.Title, &item.DatasetTitle, &item.AmountUsers, &item.AddedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
