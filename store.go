package togo

import (
	"time"
)

type Store interface {
	// Bare tasks

	AddOrUpdateTask(Task) error
	RemoveTask(Task) error
	FindTaskByName(string) (Task, error)
	FindByDueDate(*time.Time) ([]Task, error)
	OverdueTasks() ([]Task, error)
	Count() (int, error)
	All() ([]Task, error)
}

type ProjectStore interface {
	// AddOrUpdateProject will update _only_ the project. Additional steps need to be taken to update the
	// individual tasks in the project itself.
	AddOrUpdateProject(project Project) error
	TasksByPriority(string) ([]Task, error)
	AddTask(Project, Task) error
}
