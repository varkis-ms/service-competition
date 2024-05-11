package storage

import "github.com/varkis-ms/service-competition/internal/pkg/database/postgresdb"

type Storage struct {
	*postgresdb.Postgres
}

func New(db *postgresdb.Postgres) Repository {
	return &Storage{
		Postgres: db,
	}
}
