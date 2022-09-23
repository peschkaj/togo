package togo

type Store interface {
	AddOrUpdateTask(Task)
	RemoveTask(Task) bool
	FindTaskByName(string) (Task, bool)
	Count() int
}
