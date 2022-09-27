package postgres

import (
	"context"
	"fmt"
	"github.com/jaswdr/faker"
	"github.com/peschkaj/togo"
	"testing"
	"time"
)

const connectionString string = "postgres://togo_user:forhere@localhost:5432/togo"

func daysFromNow(days int) time.Time {
	durationString := fmt.Sprintf("%dh", days*24)
	duration, _ := time.ParseDuration(durationString)
	return time.Now().UTC().Add(time.Hour * duration)
}

func TestTaskCanBePersisted(t *testing.T) {
	pg := NewPgStore(connectionString)
	f := faker.New()

	taskName := f.Person().Name()

	t.Cleanup(func() {
		_ = pg.queries.RemoveTask(context.TODO(), taskName)
	})

	dueDate := daysFromNow(3)

	task := togo.Task{Name: taskName, Description: f.Lorem().Paragraph(3)}
	task.AddDueDate(dueDate)

	pg.AddOrUpdateTask(task)
}

func TestTaskCanBeRemoved(t *testing.T) {
	pg := NewPgStore(connectionString)
	f := faker.New()

	taskName := f.Person().Name()

	task := togo.Task{
		Name:        taskName,
		Description: f.Lorem().Paragraph(3),
	}
	pg.AddOrUpdateTask(task)
	result := pg.RemoveTask(task)

	if !result {
		t.Error("unable to remove task")
	}
}

func TestTaskCanBeRetrievedByName(t *testing.T) {
	pg := NewPgStore(connectionString)
	f := faker.New()

	taskName := f.Person().Name()
	expected := togo.NewTask(taskName, f.Lorem().Paragraph(3))
	pg.AddOrUpdateTask(expected)

	outcome, found := pg.FindTaskByName(taskName)

	if !found {
		t.Error("unable to find task by name")
	}

	if outcome.Name != expected.Name {
		t.Error("names do not match")
	}

	if outcome.Description != expected.Description {
		t.Error("descriptions do not match")
	}

	if outcome.Created.Equal(expected.Created) {
		t.Error("creation dates do not match")
	}

	if outcome.DueOn() != expected.DueOn() {
		t.Error("due dates do not match")
	}

	if outcome.Completed != expected.Completed {
		t.Error("completion dates do not match")
	}
}
