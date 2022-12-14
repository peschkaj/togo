-- create schema
CREATE SCHEMA togo;

CREATE TABLE IF NOT EXISTS togo.tasks (
    id BIGSERIAL PRIMARY KEY ,
    name VARCHAR(100) NOT NULL,
    description VARCHAR NOT NULL,
    created_on TIMESTAMPTZ(6) NOT NULL,
    completed_on TIMESTAMPTZ(6) NULL,
    due_date TIMESTAMPTZ(6) NULL
);

-- name: AddOrUpdateTask :exec
INSERT INTO togo.tasks (name, description, created_on, completed_on, due_date)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (name) DO UPDATE
    SET description = $2, created_on = $3, completed_on = $4, due_date = $5;

-- name: FindByName :one
SELECT * FROM togo.tasks WHERE name = $1;

-- name: FindByDueDate :many
SELECT * FROM togo.tasks WHERE due_date = $1;

-- name: FindOverdueTasks :many
SELECT * FROM togo.tasks WHERE due_date < $1;

-- name: CountTasks :one
SELECT COUNT(*) FROM togo.tasks;

-- name: AllTasks :many
SELECT * FROM togo.tasks;

-- name: RemoveTask :exec
DELETE FROM togo.tasks WHERE name = $1;
