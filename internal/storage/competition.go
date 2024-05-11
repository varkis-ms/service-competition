package storage

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/varkis-ms/service-competition/internal/model"
)

func (r *Storage) SaveCompetition(ctx context.Context, in *model.Competition) (int64, error) {
	sql, args, _ := r.Builder.
		Insert("competition").
		Columns("user_id", "title", "description", "dataset_title", "dataset_description").
		Values(in.UserID, in.Title, in.Description, in.DatasetTitle, in.DatasetDescription).
		Suffix("Returning id").
		ToSql()

	var id int64
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		if err.Error() == errDuplicate {
			return 0, model.ErrCompExists
		}
		// TODO: подумать над ситуацией, когда соревнование с таким названием уже существует
		return 0, err
	}

	return id, nil
}

func (r *Storage) EditCompetition(ctx context.Context, in *model.CompetitionEdit) error {
	// TODO: сделать транзацкии, возвращать юзер айди и сверяться от этого выдать ошибку
	sqlQuery := r.Builder.
		Update("competition").
		Where(sq.Eq{"id": in.CompetitionID}).
		Where(sq.Eq{"user_id": in.UserID})

	var isNeedUpdate bool
	if in.Title != nil {
		sqlQuery = sqlQuery.Set("title", in.Title)
		isNeedUpdate = true
	}
	if in.Description != nil {
		sqlQuery = sqlQuery.Set("description", in.Description)
		isNeedUpdate = true
	}
	if in.DatasetTitle != nil {
		sqlQuery = sqlQuery.Set("dataset_title", in.DatasetTitle)
		isNeedUpdate = true
	}
	if in.DatasetDescription != nil {
		sqlQuery = sqlQuery.Set("dataset_description", in.DatasetDescription)
		isNeedUpdate = true
	}

	if !isNeedUpdate {
		return nil
	}

	sql, args, _ := sqlQuery.ToSql()
	res, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		if err.Error() == errDuplicate {
			return model.ErrCompExists
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ErrCompNotFound
		}

		return err
	}

	if res.RowsAffected() == 0 {
		return model.ErrNoAccessToComp
	}

	return nil
}

func (r *Storage) GetCompetitionList(ctx context.Context) (model.CompetitionList, error) {
	sql, args, _ := r.Builder.
		Select("id", "title", "dataset_title").
		From("competition").
		ToSql()

	var competitionList model.CompetitionList
	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var competition model.Competition

		err = rows.Scan(&competition.CompetitionID, &competition.Title, &competition.DatasetTitle)
		if err != nil {
			return nil, err
		}

		competitionList = append(competitionList, competition)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return competitionList, nil
}

func (r *Storage) GetCompetitionInfo(ctx context.Context, competitionID int64) (*model.Competition, error) {
	sql, args, _ := r.Builder.
		Select("id", "title", "description", "dataset_title", "dataset_description").
		From("competition").
		Where(sq.Eq{"id": competitionID}).
		ToSql()

	var comp model.Competition
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&comp.CompetitionID,
		&comp.Title,
		&comp.Description,
		&comp.DatasetTitle,
		&comp.DatasetDescription,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrCompNotFound
		}

		return nil, err
	}

	return &comp, nil
}
