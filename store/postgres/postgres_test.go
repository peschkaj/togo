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
