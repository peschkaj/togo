package store

import (
	"github.com/jaswdr/faker"
	"github.com/peschkaj/togo"
	"testing"
)

func TestMultipleCallsToAddIncreaseCount(t *testing.T) {
	ms := NewMemoryStore()
	f := faker.New()

	initialCount := ms.Count()
	if initialCount != 0 {
		t.Error("memory store is not empty after initialization")
	}

	var previousCount = initialCount

	for i := 0; i < 3; i++ {
		ms.AddOrUpdateTask(togo.NewTask(f.Person().Name(), f.Lorem().Paragraph(3)))
		count := ms.Count()

		if count <= previousCount {
			t.Error("count did not increment after AddOrUpdateTask()")
		}
	}
}

func TestTasksCanBeRetrievedByName(t *testing.T) {
	ms := NewMemoryStore()
	f := faker.New()

	task := togo.NewTask(f.Person().Name(), f.Lorem().Paragraph(1))

	// add the task we want to find
	ms.AddOrUpdateTask(task)
	// and several other tasks to ensure we're not testing a degenerate case
	for i := 0; i < 3; i++ {
		ms.AddOrUpdateTask(togo.NewTask(f.Person().Name(), f.Lorem().Paragraph(3)))
	}

	otherTask, found := ms.FindTaskByName(task.Name)
	if !found {
		t.Error("unable to find task by name")
	}

	if task != *otherTask {
		t.Error("found task wasn't the same as original task")
	}
}

func TestRemovedTaskCannotBeFound(t *testing.T) {
	ms := NewMemoryStore()
	f := faker.New()
	originalCount := ms.Count()

	task := togo.NewTask(f.Person().Name(), f.Lorem().Paragraph(1))
	ms.AddOrUpdateTask(task)
	for i := 0; i < 3; i++ {
		ms.AddOrUpdateTask(togo.NewTask(f.Person().Name(), f.Lorem().Paragraph(3)))
	}

	if !(ms.Count() > originalCount) {
		t.Error("current count is not greater than original count")
	}

	removed := ms.RemoveTask(task)
	if !removed {
		t.Error("unable to remove task")
	}

	_, found := ms.FindTaskByName(task.Name)
	if found {
		t.Error("found task in store but should not be able to")
	}
}
