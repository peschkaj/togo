package memory

import (
	"fmt"
	"github.com/jaswdr/faker"
	"github.com/peschkaj/togo"
	"testing"
	"time"
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

func TestUpdatedTaskIsChanged(t *testing.T) {
	ms := NewMemoryStore()
	f := faker.New()

	task := togo.NewTask(f.Person().Name(), f.Lorem().Paragraph(1))
	ms.AddOrUpdateTask(task)
	originalTask := task

	task.Description = f.Lorem().Paragraph(3)
	ms.AddOrUpdateTask(task)

	updatedTask, found := ms.FindTaskByName(task.Name)
	if !found {
		t.Error("task saved but not found after update")
	}

	if updatedTask.Description == originalTask.Description {
		t.Error("task was not successfully updated")
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

func TestFindByNameReturnsFalseWhenNotFound(t *testing.T) {
	ms := NewMemoryStore()
	f := faker.New()

	// add several tasks to ensure we're not testing a degenerate case
	for i := 0; i < 3; i++ {
		ms.AddOrUpdateTask(togo.NewTask(f.Person().Name(), f.Lorem().Paragraph(3)))
	}

	_, found := ms.FindTaskByName("asdf")

	if found != false {
		t.Error("searched for a task that does not exist and found it")
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

func TestOverdueTasksCanBeRetrieved(t *testing.T) {
	ms := NewMemoryStore()
	f := faker.New()
	// start with two days ago
	start := -2

	for i := 0; i < 10; i++ {
		start += i
		task := togo.NewTask(f.Person().Name(), f.Lorem().Paragraph(3))
		durationString := fmt.Sprintf("%dh", start*24)
		duration, _ := time.ParseDuration(durationString)
		dueDate := time.Now().Add(time.Hour * duration)
		task.AddDueDate(dueDate)

		ms.AddOrUpdateTask(task)
	}

	overdueTasks := ms.OverdueTasks()
	if overdueTasks == nil {
		t.Error("no overdue tasks found")
	}

	if len(overdueTasks) != 5 {
		t.Error("not all overdue tasks retrieved")
	}
}

func TestTasksCanBeRetrievedByDueDate(t *testing.T) {
	ms := NewMemoryStore()
	f := faker.New()
	// start with two days ago
	start := -2

	now := time.Now()

	for i := 0; i < 10; i++ {
		start += i
		if i < 5 {
			start = 0
		}
		task := togo.NewTask(f.Person().Name(), f.Lorem().Paragraph(3))
		durationString := fmt.Sprintf("%dh", start*24)
		duration, _ := time.ParseDuration(durationString)
		dueDate := now.Add(time.Hour * duration)
		task.AddDueDate(dueDate)

		ms.AddOrUpdateTask(task)
	}

	tasks := ms.FindByDueDate(&now)
	if tasks == nil {
		t.Error("no tasks found")
	}

	if len(tasks) != 5 {
		t.Error("all tasks not found")
	}
}

func TestTasksCanBeRetrievedByNilDueDate(t *testing.T) {
	ms := NewMemoryStore()
	f := faker.New()

	task := togo.NewTask(f.Person().Name(), f.Lorem().Paragraph(1))
	ms.AddOrUpdateTask(task)

	tasks := ms.FindByDueDate(nil)
	if tasks == nil {
		t.Error("no tasks found with nil due date")
	}

	if len(tasks) != 1 {
		t.Error("incorrect number of tasks found with nil due date")
	}
}
