package store

import (
	"github.com/peschkaj/togo"
	"time"
)

type Store interface {
	AddOrUpdateTask(togo.Task)
	RemoveTask(togo.Task) bool
	FindTaskByName(string) (*togo.Task, bool)
	FindByDueDate(*time.Time) []togo.Task
	OverdueTasks() []togo.Task
	Count() int
	All() []togo.Task
}
