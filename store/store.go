package store

import (
	"github.com/peschkaj/togo"
	"time"
)

type Store interface {
	AddOrUpdateTask(togo.Task) error
	RemoveTask(togo.Task) error
	FindTaskByName(string) (togo.Task, error)
	FindByDueDate(*time.Time) ([]togo.Task, error)
	OverdueTasks() ([]togo.Task, error)
	Count() (int, error)
	All() ([]togo.Task, error)
}
