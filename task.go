package togo

import (
	"time"
)

type Task struct {
	Name        string
	Description string
	Created     time.Time
	Completed   *time.Time
	DueDate     *time.Time
}

func NewTask(name, description string) Task {
	return Task{Name: name, Description: description, Created: time.Now()}
}

func (t *Task) IsCompleted() bool {
	return t.Completed != nil && t.Completed.After(time.Now())
}

func (t *Task) Overdue() bool {
	return t.DueDate != nil && t.DueDate.Before(time.Now())
}

func (t *Task) Complete() {
	completionTime := time.Now()
	t.Completed = &completionTime
}

func (t *Task) CompletionDate() *time.Time {
	return t.Completed
}

// AddDueDate strips off the time component and stores the result as the due date
func (t *Task) AddDueDate(due time.Time) {
	yyyy, mm, dd := due.Date()
	newDate := time.Date(yyyy, mm, dd, 0, 0, 0, 0, time.UTC)
	t.DueDate = &newDate
}

func (t *Task) DueOn() *time.Time {
	return t.DueDate
}
