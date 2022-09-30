package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/peschkaj/togo"
	"time"
)

type PgStore struct {
	pool *pgxpool.Pool
}

const addOrUpdateTask = `-- name: AddOrUpdateTask 
INSERT INTO togo.tasks (name, description, created_on, completed_on, due_date)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (name) DO UPDATE
    SET description = $2, created_on = $3, completed_on = $4, due_date = $5;
`

const removeTask = `-- name: RemoveTask
DELETE FROM togo.tasks WHERE name = $1;
`

const findTaskByName = `-- name: FindTaskByName 
SELECT name, description, created_on as created, completed_on as completed, due_date
FROM togo.tasks 
WHERE name = $1;
`

const findTasksByDueDate = `-- name: FindTasksByDueDate
SELECT name, description, created_on as created, completed_on as completed, due_date
FROM togo.tasks 
WHERE due_date BETWEEN $1 AND $2;
`

const findOverdueTasks = `-- name: FindOverdueTasks
SELECT name, description, created_on as created, completed_on as completed, due_date 
FROM togo.tasks 
WHERE due_date < CURRENT_TIMESTAMP;
`

func NewPgStore(connectionURI string) PgStore {
	p, err := pgxpool.New(context.TODO(), connectionURI)
	if err != nil {
		panic("cannot connect to postgres backing store")
	}

	return PgStore{pool: p}
}

func (p PgStore) AddOrUpdateTask(t togo.Task) error {
	_, err := p.pool.Exec(context.TODO(),
		addOrUpdateTask,
		t.Name,
		t.Description,
		t.Created,
		t.Completed,
		t.DueOn(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (p PgStore) RemoveTask(taskName string) error {
	_, err := p.pool.Exec(context.TODO(),
		removeTask,
		taskName)
	if err != nil {
		return errors.New("unable to remove task")
	}
	return nil
}

func (p PgStore) FindTaskByName(name string) (togo.Task, error) {
	row := p.pool.QueryRow(context.TODO(), findTaskByName, name)
	var i togo.Task
	err := row.Scan(
		&i.Name,
		&i.Description,
		&i.Created,
		&i.Completed,
		&i.DueDate,
	)

	return i, err
}

func (p PgStore) FindTasksByDueDate(d time.Time) ([]togo.Task, error) {
	start := timeToDate(d)
	end := timeToDate(d).Add(24 * time.Hour)

	rows, err := p.pool.Query(context.TODO(), findTasksByDueDate, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tasks := []togo.Task{}
	for rows.Next() {
		var t togo.Task
		if err := rows.Scan(
			&t.Name,
			&t.Description,
			&t.Created,
			&t.Completed,
			&t.DueDate,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (p PgStore) FindOverdueTasks() ([]togo.Task, error) {
	rows, err := p.pool.Query(context.TODO(), findOverdueTasks)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	tasks := []togo.Task{}
	for rows.Next() {
		var t togo.Task
		if err := rows.Scan(
			&t.Name,
			&t.Description,
			&t.Created,
			&t.Completed,
			&t.DueDate,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func timeToDate(t time.Time) time.Time {
	yyyy, mm, dd := t.Date()
	return time.Date(yyyy, mm, dd, 0, 0, 0, 0, t.Location())
}
