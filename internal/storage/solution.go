package storage

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

const (
	statusInProgress = "in progress"
	statusInQueue    = "in queue"
	statusDone       = "done"
)

func (r *Storage) SaveSolutionInfo(ctx context.Context, solutionID int64, score float32, interval time.Duration) error {
	sql, args, _ := r.Builder.
		Update("leaderboard").
		Set("status", statusDone).
		Set("score", score).
		Set("run_time", interval).
		Where(sq.Eq{"id": solutionID}).
		ToSql()

	_, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *Storage) GetNextSolutionInfo(ctx context.Context) (int64, int64, int64, error) {
	subQuery, subArgs, _ := r.Builder.
		Select("id").
		From("leaderboard").
		Where(sq.Eq{"status": statusInQueue}).
		OrderBy("added_at").
		Limit(1).
		ToSql()

	var id int64
	err := r.Pool.QueryRow(ctx, subQuery, subArgs...).Scan(
		&id,
	)
	if err != nil {
		return 0, 0, 0, err
	}

	sql, args, _ := r.Builder.
		Update("leaderboard").
		Set("status", statusInProgress).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING user_id, competition_id, id").
		ToSql()

	var userID, compID, solutionID int64
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&userID,
		&compID,
		&solutionID,
	)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return 0, 0, 0, err
		}
	}

	return userID, compID, solutionID, nil
}
