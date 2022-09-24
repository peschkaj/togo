package memory

import (
	"github.com/peschkaj/togo"
	art "github.com/plar/go-adaptive-radix-tree"
	"time"
)

type InMemoryStore struct {
	ts        art.Tree
	byDueDate art.Tree
}

func NewMemoryStore() InMemoryStore {
	return InMemoryStore{ts: art.New(), byDueDate: art.New()}
}

func (ms InMemoryStore) AddOrUpdateTask(t togo.Task) {
	ms.ts.Insert(art.Key(t.Name), t)
	addOrUpdateByDueDate(ms.byDueDate, t)
}

func (ms InMemoryStore) RemoveTask(t togo.Task) bool {
	_, removed := ms.ts.Delete(art.Key(t.Name))
	removeByDueDate(ms.byDueDate, t)
	return removed
}

func (ms InMemoryStore) FindTaskByName(name string) (*togo.Task, bool) {
	value, found := ms.ts.Search(art.Key(name))

	if !found {
		return nil, found
	}

	switch t := value.(type) {
	case togo.Task:
		return &t, true
	default:
		return nil, false
	}
}

func (ms InMemoryStore) FindByDueDate(dueDate *time.Time) []togo.Task {
	key := dateToKey(dueDate)
	value, found := ms.byDueDate.Search(key)
	if !found {
		return nil
	}

	switch tasks := value.(type) {
	case []togo.Task:
		return tasks
	default:
		panic("type mismatch reading from index")
	}
}

func (ms InMemoryStore) Count() int {
	return ms.ts.Size()
}

func (ms InMemoryStore) All() []togo.Task {
	items := []togo.Task{}

	iter := ms.ts.Iterator()

	for iter.HasNext() {
		node, err := iter.Next()
		if err != nil {
			break
		}

		switch t := node.Value().(type) {
		case togo.Task:
			items = append(items, t)
		default:
			continue
		}
	}

	return items
}

func (ms InMemoryStore) OverdueTasks() []togo.Task {
	iter := ms.byDueDate.Iterator()
	var tasks []togo.Task
	now := time.Now().UTC()

loop:
	for iter.HasNext() {
		node, err := iter.Next()
		if err != nil {
			return tasks
		}

		switch nodeTasks := node.Value().(type) {
		case []togo.Task:
			if len(nodeTasks) == 0 {
				continue
			}

			if nodeTasks[0].DueDate() == nil {
				continue
			}

			if nodeTasks[0].DueDate().After(now) {
				break loop
			}
			tasks = append(tasks, nodeTasks[:]...)
		}
	}

	return tasks
}

func addOrUpdateByDueDate(tree art.Tree, t togo.Task) {
	key := dateToKey(t.DueDate())
	updateIndex(tree, key, t)
}

func removeByDueDate(tree art.Tree, t togo.Task) bool {
	key := dateToKey(t.DueDate())

	return removeFromIndex(tree, key, t)
}

func dateToKey(date *time.Time) art.Key {
	var key art.Key = []byte{}
	if date != nil {
		yyyy, mm, dd := date.Date()
		newDate := time.Date(yyyy, mm, dd, 0, 0, 0, 0, time.UTC)
		key, _ = newDate.GobEncode()
	}
	return key
}

func removeFromIndex(tree art.Tree, key art.Key, t togo.Task) bool {
	value, found := tree.Search(key)
	// key not found, don't need to delete
	if !found {
		return false
	}

	switch tasks := value.(type) {
	case []togo.Task:
		originalLength := len(tasks)
		// delete an element from the array
		newLength := 0
		for i := range tasks {
			if t.Name != tasks[i].Name {
				tasks[newLength] = tasks[i]
				newLength++
			}
		}

		// re-slice the array to remove extra index
		tasks = tasks[:newLength]

		// didn't find it
		if newLength == originalLength {
			return false
		}

		// this was the last item for this index, remove it
		if newLength == 0 {
			tree.Delete(key)
			return true
		}

		// wasn't the last item, need to update the indexed values
		tree.Insert(key, tasks)
		return true
	}

	return false
}

func updateIndex(tree art.Tree, key art.Key, t togo.Task) {
	value, found := tree.Search(key)
	// key not found
	if !found {
		tree.Insert(key, []togo.Task{t})
		return
	}

	// update by searching for the task by name
	switch tasks := value.(type) {
	case []togo.Task:
		for i, task := range tasks {
			if task.Name == t.Name {
				// found it, update the list in place
				tasks[i] = t
				// save back to the tree and bail
				_, _ = tree.Insert(key, tasks)
				return
			}
		}

		// didn't find the task, eh?
		tasks = append(tasks, t)
		tree.Insert(key, tasks)
		return
	default:
		panic("type mismatch in index")
	}
}
