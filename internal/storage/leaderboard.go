package storage

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/varkis-ms/service-competition/internal/model"
)

func (r *Storage) GetLeaderboard(ctx context.Context, competitionID int64) (model.LeaderBoardList, error) {
	subQuery := fmt.Sprintf(
		"(SELECT user_id, score, added_at, ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY score DESC) "+
			"AS row_num FROM leaderboard WHERE competition_id = %d) AS ranked", competitionID,
	)

	sql, args, _ := r.Builder.
		Select("user_id", "score", "added_at").
		From(subQuery).
		Where("row_num = 1").
		OrderBy("score DESC").
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}
	defer rows.Close()

	var leaderboardList model.LeaderBoardList
	for rows.Next() {
		var leaderboard model.LeaderBoard

		err = rows.Scan(&leaderboard.UserID, &leaderboard.Score, &leaderboard.AddedAt)
		if err != nil {
			return nil, err
		}

		leaderboardList = append(leaderboardList, leaderboard)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return leaderboardList, nil
}

func (r *Storage) GetAmountUsersAndMaxScore(ctx context.Context, competitionID int64) (float32, int64, error) {
	sql, args, _ := r.Builder.
		Select("MAX(score)", "COUNT(DISTINCT user_id)").
		From("leaderboard").
		Where(sq.Eq{"competition_id": competitionID}).
		ToSql()

	var maxScore float32
	var amountUsers int64
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&maxScore,
		&amountUsers,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, 0, nil
		}

		return 0, 0, err
	}

	return maxScore, amountUsers, nil
}

func (r *Storage) SaveScore(ctx context.Context, userID, competitionID, score int64) error {
	// TODO: runtime
	sql, args, _ := r.Builder.
		Update("leaderboard").
		Set("score", score).
		Set("queue_id", nil).
		Where(sq.Eq{"user_id": userID}).
		Where(sq.Eq{"competition_id": competitionID}).
		ToSql()

	_, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
