package postgres

import (
	"context"
	"database/sql"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/peschkaj/togo"
	"time"
)

type PgStore struct {
	pool    *pgxpool.Pool
	queries *Queries
}

const addOrUpdateTask = `-- name: AddOrUpdateTask :exec
INSERT INTO togo.tasks (name, description, created_on, completed_on, due_date)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (name) DO UPDATE
    SET description = $2, created_on = $3, completed_on = $4, due_date = $5
`

func NewPgStore(connectionURI string) PgStore {
	p, err := pgxpool.New(context.TODO(), connectionURI)
	if err != nil {
		panic("cannot connect to postgres backing store")
	}

	q := New(p)
	return PgStore{pool: p, queries: q}
}

func (p PgStore) AddOrUpdateTask(t togo.Task) error {
	_, err := p.pool.Exec(
		context.TODO(),
		addOrUpdateTask,
		t.Name,
		t.Description,
		t.Created,
		t.Completed,
		t.DueOn())

	return err
}

func (p PgStore) RemoveTask(taskName string) error {
	_, err := p.pool.Exec(
		context.TODO(),
		removeTask,
		taskName)

	return err
}

const findByName = `-- name: FindByName :one
SELECT name, description, created_on::timestamptz, completed_on::timestamptz, due_date::timestamptz FROM togo.tasks WHERE name = $1
`

func (p PgStore) FindTaskByName(name string) (togo.Task, error) {
	var t togo.Task

	rows, err := p.pool.Query(context.TODO(),
		findByName,
		name)
	defer rows.Close()
	if err != nil {
		return t, err
	}

	if rows.Next() {
		var taskName, description string
		var created pgtype.Timestamptz
		var completed pgtype.Timestamptz
		var dueDate pgtype.Timestamptz
		err := rows.Scan(&taskName, &description, &created, &completed, &dueDate)
		if err != nil {
			return t, err
		}

		t.Name = taskName
		t.Description = description
		t.Created = created.Time
		t.Completed = timestamptzToTime(completed)

		d := timestamptzToTime(dueDate)
		if d != nil {
			t.AddDueDate(*d)
		}
	}

	return t, nil
}

func timestamptzToTime(ts pgtype.Timestamptz) *time.Time {
	if ts.Status == pgtype.Null {
		return nil
	}
	t := ts.Time
	return &t
}

//func (p PgStore) FindTaskByName(name string) (*togo.Task, bool) {
//	byName, err := p.queries.FindByName(context.TODO(), name)
//	if err != nil {
//		return nil, false
//	}
//
//	var completed *time.Time
//	if byName.CompletedOn.Status == pgtype.Present {
//		completed = &byName.CompletedOn.Time
//	}
//
//	task := togo.Task{
//		Name:        byName.Name,
//		Description: byName.Description,
//		Created:     byName.CreatedOn.Time,
//		Completed:   completed,
//	}
//
//	if byName.DueDate.Status == pgtype.Present {
//		task.AddDueDate(byName.DueDate.Time)
//	}
//
//	return &task, true
//}

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
