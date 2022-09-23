package togo

import (
	"time"
)

type Task struct {
	Name        string
	description string
	created     time.Time
	completed   *time.Time
	dueDate     *time.Time
}

func NewTask(name, description string) Task {
	return Task{Name: name, description: description, created: time.Now().UTC()}
}

func (t *Task) Completed() bool {
	return t.completed != nil && t.completed.After(time.Now().UTC())
}

func (t *Task) Overdue() bool {
	return t.dueDate != nil && t.dueDate.Before(time.Now().UTC())
}

func (t *Task) Complete() {
	completionTime := time.Now().UTC()
	t.completed = &completionTime
}

func (t *Task) CompletionDate() *time.Time {
	return t.completed
}

func (t *Task) AddDueDate(due time.Time) {
	newDate := due
	t.dueDate = &newDate
}
