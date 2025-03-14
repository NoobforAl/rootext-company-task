package repository

import (
	"context"
	"ratblog/contract"
	"ratblog/internal/database"

	"github.com/jackc/pgx/v5"
)

type repository struct {
	db *database.Queries
}

func New(ctx context.Context, postgresUrl string) contract.Repository {
	conn, err := pgx.Connect(ctx, postgresUrl)
	if err != nil {
		panic(err)
	}

	return &repository{
		db: database.New(conn),
	}
}
