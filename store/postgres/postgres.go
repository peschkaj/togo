package postgres

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/peschkaj/togo"
	"time"
)

type PgStore struct {
	pool    *pgxpool.Pool
	queries *Queries
}

func NewPgStore(connectionURI string) PgStore {
	pool, err := pgxpool.Connect(context.Background(), connectionURI)
	if err != nil {
		panic("cannot connect to postgres backing store")
	}
	q := New(pool)
	return PgStore{pool: pool, queries: q}
}

func (p PgStore) AddOrUpdateTask(t togo.Task) {
	due := timeToNullTime(t.DueDate())
	completed := timeToNullTime(t.CompletionDate())
	params := AddOrUpdateTaskParams{
		Name:        t.Name,
		Description: t.Description,
		CreatedOn:   t.Created(),
		DueDate:     due,
		CompletedOn: completed,
	}

	err := p.queries.AddOrUpdateTask(context.TODO(), params)
	if err != nil {
		panic("unable to add or update task")
	}
}

func timeToNullTime(time *time.Time) sql.NullTime {
	if time == nil {
		return sql.NullTime{Valid: false}
	}

	return sql.NullTime{Time: *time, Valid: true}
}
