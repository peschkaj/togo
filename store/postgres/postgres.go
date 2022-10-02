package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/peschkaj/togo"
	"time"
)

type PgStore struct {
	pool *pgxpool.Pool
}

const addOrUpdateTask = `-- name: AddOrUpdateTask 
INSERT INTO togo.tasks (name, description, priority, project_id, created_on, completed_on, due_date)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (name) DO UPDATE
    SET description = $2, priority = $3, project_id = $4, created_on = $5, completed_on = $6, due_date = $7;
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

const countTasks = `-- name: CountTasks
SELECT COUNT(*) FROM togo.Tasks;
`

const allTasks = `-- name: AllTasks
SELECT name, description, created_on as created, completed_on as completed, due_date 
FROM togo.tasks 
`

func NewPgStore(connectionURI string) PgStore {
	p, err := pgxpool.New(context.TODO(), connectionURI)
	if err != nil {
		panic("cannot connect to postgres backing store")
	}

	return PgStore{pool: p}
}

func (p PgStore) AddOrUpdateTask(t togo.Task, projectId *int64) error {
	_, err := p.pool.Exec(context.TODO(),
		addOrUpdateTask,
		t.Name,
		t.Description,
		t.Priority,
		projectId,
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

func (p PgStore) Count() (int, error) {
	row := p.pool.QueryRow(context.TODO(), countTasks)
	var count int
	err := row.Scan(&count)

	return count, err
}

func (p PgStore) All() ([]togo.Task, error) {
	rows, err := p.pool.Query(context.TODO(), allTasks)
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

const addOrUpdateProject = `-- name: AddOrUpdateProject
INSERT INTO togo.projects (name, description) 
VALUES ($1, $2)
ON CONFLICT (name) DO UPDATE
	SET description = $2;
`

const tasksByPriority = `-- name: TasksByPriority
SELECT t.name, t.description, t.priority, t.created_on, t.completed_on, t.due_date
FROM togo.tasks AS t
	JOIN togo.projects AS p ON t.project_id = p.id
WHERE p.name = $1
ORDER BY t.priority DESC, t.due_date ASC;
`

const findProjectIdByName = `-- name: FindProjectId
SELECT id FROM togo.projects WHERE name = $1;
`

func (p PgStore) AddOrUpdateProject(project togo.Project) error {
	_, err := p.pool.Exec(context.TODO(), addOrUpdateProject, project.Name, project.Description)
	return err
}

func (p PgStore) TasksByPriority(projectName string) ([]togo.Task, error) {
	rows, err := p.pool.Query(context.TODO(), tasksByPriority, projectName)
	if err != nil {
		return nil, err
	}

	tasks := []togo.Task{}
	for rows.Next() {
		var t togo.Task
		if err := rows.Scan(
			&t.Name,
			&t.Description,
			&t.Priority,
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

func (p PgStore) AddTaskToProject(project togo.Project, t togo.Task) error {
	// look up the project's ID by name
	projectId, err := p.findProjectId(project.Name)
	if err != nil {
		return err
	}

	// AddOrUpdateTask :)
	if err := p.AddOrUpdateTask(t, &projectId); err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("project not found")
		}
		return err
	}

	return nil
}

func (p PgStore) findProjectId(projectName string) (int64, error) {
	row := p.pool.QueryRow(context.TODO(), findProjectIdByName, projectName)
	var projectId int64
	if err := row.Scan(&projectId); err != nil {
		return -1, err
	}

	return projectId, nil
}
