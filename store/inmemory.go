package store

import (
	"github.com/peschkaj/togo"
	art "github.com/plar/go-adaptive-radix-tree"
)

type MemoryStore struct {
	ts art.Tree
}

func NewMemoryStore() MemoryStore {
	return MemoryStore{ts: art.New()}
}

func (ms *MemoryStore) AddOrUpdateTask(t togo.Task) {
	ms.ts.Insert(art.Key(t.Name), t)
}

func (ms *MemoryStore) RemoveTask(t togo.Task) bool {
	_, removed := ms.ts.Delete(art.Key(t.Name))
	return removed
}

func (ms *MemoryStore) FindTaskByName(name string) (*togo.Task, bool) {
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

func (ms *MemoryStore) Count() int {
	return ms.ts.Size()
}

func (ms *MemoryStore) All() []togo.Task {
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
