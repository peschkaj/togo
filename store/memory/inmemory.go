package memory

import (
	"github.com/peschkaj/togo"
	art "github.com/plar/go-adaptive-radix-tree"
	"sort"
	"time"
)

type InMemoryStore struct {
	ts        art.Tree
	byDueDate art.Tree
}

func NewMemoryStore() InMemoryStore {
	return InMemoryStore{ts: art.New(), byDueDate: art.New()}
}

func SortByPriority(ts togo.Tasks) []togo.Task {
	if len(ts) == 0 {
		return ts
	}

	sorted := ts
	sort.Sort(sorted)
	return sorted
}

func (ms InMemoryStore) AddOrUpdateTask(t togo.Task) error {
	ms.ts.Insert(art.Key(t.Name), t)
	addOrUpdateByDueDate(ms.byDueDate, t)
	return nil
}

func (ms InMemoryStore) RemoveTask(t togo.Task) error {
	ms.ts.Delete(art.Key(t.Name))
	removeByDueDate(ms.byDueDate, t)
	return nil
}

func (ms InMemoryStore) FindTaskByName(name string) (togo.Task, error) {
	value, found := ms.ts.Search(art.Key(name))

	if !found {
		return togo.Task{}, nil
	}

	switch t := value.(type) {
	case togo.Task:
		return t, nil
	default:
		panic("type mismatch in index")
	}
}

func (ms InMemoryStore) FindByDueDate(dueDate *time.Time) ([]togo.Task, error) {
	key := dateToKey(dueDate)
	value, found := ms.byDueDate.Search(key)
	if !found {
		return nil, nil
	}

	switch tasks := value.(type) {
	case []togo.Task:
		return tasks, nil
	default:
		panic("type mismatch reading from index")
	}
}

func (ms InMemoryStore) Count() int {
	return ms.ts.Size()
}

func (ms InMemoryStore) All() ([]togo.Task, error) {
	items := []togo.Task{}

	iter := ms.ts.Iterator()

	for iter.HasNext() {
		node, err := iter.Next()
		if err != nil {
			return nil, err
		}

		switch t := node.Value().(type) {
		case togo.Task:
			items = append(items, t)
		default:
			continue
		}
	}

	return items, nil
}

func (ms InMemoryStore) OverdueTasks() ([]togo.Task, error) {
	iter := ms.byDueDate.Iterator()
	var tasks []togo.Task
	now := time.Now()

loop:
	for iter.HasNext() {
		node, err := iter.Next()
		if err != nil {
			return tasks, nil
		}

		switch nodeTasks := node.Value().(type) {
		case []togo.Task:
			if len(nodeTasks) == 0 {
				continue
			}

			if nodeTasks[0].DueOn() == nil {
				continue
			}

			if nodeTasks[0].DueOn().After(now) {
				break loop
			}
			tasks = append(tasks, nodeTasks[:]...)
		}
	}

	return tasks, nil
}

func addOrUpdateByDueDate(tree art.Tree, t togo.Task) {
	key := dateToKey(t.DueOn())
	updateIndex(tree, key, t)
}

func removeByDueDate(tree art.Tree, t togo.Task) bool {
	key := dateToKey(t.DueOn())

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
