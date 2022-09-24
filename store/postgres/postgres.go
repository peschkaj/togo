package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgStore struct {
	pool *pgxpool.Pool
}

func NewPgStore(connectionURI string) PgStore {
	pool, err := pgxpool.Connect(context.Background(), connectionURI)
	if err != nil {
		panic("cannot connect to postgres backing store")
	}

	return PgStore{pool: pool}
}
