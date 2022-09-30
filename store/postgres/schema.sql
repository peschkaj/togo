CREATE ROLE togo_user
    LOGIN
    NOSUPERUSER
    PASSWORD 'forhere';

CREATE SCHEMA IF NOT EXISTS togo;

CREATE TABLE IF NOT EXISTS togo.tasks (
    id BIGSERIAL PRIMARY KEY ,
    name VARCHAR(100) NOT NULL,
    description VARCHAR NOT NULL,
    created_on TIMESTAMPTZ(6) NOT NULL,
    completed_on TIMESTAMPTZ(6) NULL,
    due_date TIMESTAMPTZ(6) NULL
);

CREATE UNIQUE INDEX ux_tasks_name ON togo.tasks(name);
CREATE INDEX ix_tasks_due_date ON togo.tasks(due_date);

GRANT USAGE ON SCHEMA togo TO togo_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA togo TO togo_user;
GRANT SELECT, USAGE ON ALL SEQUENCES IN SCHEMA togo TO togo_user;
