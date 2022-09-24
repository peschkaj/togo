CREATE SCHEMA IF NOT EXISTS togo;

CREATE TABLE IF NOT EXISTS togo.tasks (
    id BIGSERIAL PRIMARY KEY ,
    name VARCHAR(100) NOT NULL,
    description VARCHAR NOT NULL,
    created_on timestamptz NOT NULL,
    completed_on timestamptz NULL,
    due_date timestamptz NULL
);

CREATE UNIQUE INDEX ux_tasks_name ON togo.tasks(name);
CREATE INDEX ix_tasks_due_date ON togo.tasks(due_date);