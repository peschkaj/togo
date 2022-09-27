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

	params := AddOrUpdateTaskParams{
		Name:        t.Name,
		Description: t.Description,
		CreatedOn:   t.Created,
		DueDate:     timeToNullTime(t.DueOn()),
		CompletedOn: timeToNullTime(t.Completed),
	}

	err := p.queries.AddOrUpdateTask(context.TODO(), params)
	if err != nil {
		panic("unable to add or update task")
	}
}

func (p PgStore) RemoveTask(task togo.Task) bool {
	err := p.queries.RemoveTask(context.TODO(), task.Name)
	if err != nil {
		return false
	}

	return true
}

func (p PgStore) FindTaskByName(name string) (*togo.Task, bool) {
	byName, err := p.queries.FindByName(context.TODO(), name)
	if err != nil {
		return nil, false
	}

	task := togo.Task{
		Name:        byName.Name,
		Description: byName.Description,
		Created:     byName.CreatedOn,
		Completed:   nullTimeToTime(byName.CompletedOn),
	}

	if byName.DueDate.Valid {
		task.AddDueDate(byName.DueDate.Time)
	}

	return &task, true
}

func nullTimeToTime(time sql.NullTime) *time.Time {
	if !time.Valid {
		return nil
	}

	return &time.Time
}

func timeToNullTime(time *time.Time) sql.NullTime {
	if time == nil {
		return sql.NullTime{Valid: false}
	}

	return sql.NullTime{Time: *time, Valid: true}
}
