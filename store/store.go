package store

import "github.com/peschkaj/togo"

type Store interface {
	AddOrUpdateTask(togo.Task)
	RemoveTask(togo.Task) bool
	FindTaskByName(string) (*togo.Task, bool)
	Count() int
	All() []togo.Task
}
