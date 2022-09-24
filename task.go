package togo

import (
	"time"
)

type Task struct {
	Name        string
	Description string
	created     time.Time
	completed   *time.Time
	dueDate     *time.Time
}

func NewTask(name, description string) Task {
	return Task{Name: name, Description: description, created: time.Now().UTC()}
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

// AddDueDate strips off the time component and stores the result as the due date
func (t *Task) AddDueDate(due time.Time) {
	yyyy, mm, dd := due.Date()
	newDate := time.Date(yyyy, mm, dd, 0, 0, 0, 0, time.UTC)
	t.dueDate = &newDate
}

func (t *Task) DueDate() *time.Time {
	return t.dueDate
}

func (t *Task) Created() time.Time {
	return t.created
}
