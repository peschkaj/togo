package postgres

import (
	"errors"
	"fmt"
	"github.com/jaswdr/faker"
	"github.com/peschkaj/togo"
	"testing"
	"time"
)

const connectionString string = "postgres://togo_user:forhere@localhost:5432/togo"

func daysFromNow(days int) *time.Time {
	durationString := fmt.Sprintf("%dh", days*24)
	duration, _ := time.ParseDuration(durationString)
	theTime := time.Now().UTC().Add(time.Hour * duration)
	return &theTime
}

func TestTaskCanBePersisted(t *testing.T) {
	pg := NewPgStore(connectionString)
	f := faker.New()

	taskName := f.Person().Name()

	t.Cleanup(func() {
		err := pg.RemoveTask(taskName)
		if err != nil {
			return
		}
	})

	dueDate := daysFromNow(3)

	task := togo.Task{Name: taskName, Description: f.Lorem().Paragraph(3)}
	task.AddDueDate(*dueDate)

	err := pg.AddOrUpdateTask(task)
	if err != nil {
		t.Error(err)
	}
}

func TestTaskCanBeRemoved(t *testing.T) {
	pg := NewPgStore(connectionString)
	f := faker.New()

	taskName := f.Person().Name()

	task := togo.Task{
		Name:        taskName,
		Description: f.Lorem().Paragraph(3),
	}

	if err := pg.AddOrUpdateTask(task); err != nil {
		t.Error(err)
	}

	if err := pg.RemoveTask(task.Name); err != nil {
		t.Error("unable to remove task")
	}
}

func TestSimpleTaskCanBeRetrievedByName(t *testing.T) {
	pg := NewPgStore(connectionString)
	f := faker.New()

	taskName := f.Person().Name()
	expected := togo.NewTask(taskName, f.Lorem().Paragraph(3))
	err := pg.AddOrUpdateTask(expected)
	if err != nil {
		t.Error(err)
	}

	outcome, err := pg.FindTaskByName(taskName)

	if err != nil {
		t.Error("unable to find task by name")
	}

	if outcome.Name != expected.Name {
		t.Error("names do not match")
	}

	if outcome.Description != expected.Description {
		t.Error("descriptions do not match")
	}

	if !compareTime(&expected.Created, &outcome.Created) {
		t.Error("creation dates do not match")
	}

	if !compareTime(expected.DueOn(), outcome.DueOn()) {
		t.Error("due dates do not match")
	}

	if !compareTime(expected.Completed, outcome.Completed) {
		t.Error("completion times do not match")
	}
}

func TestTasksCanBeRetrievedByName(t *testing.T) {
	pg := NewPgStore(connectionString)
	f := faker.New()

	testCases := []struct {
		task    togo.Task
		dueDate *time.Time
	}{
		{task: togo.Task{Name: f.Person().Name(), Description: f.Lorem().Paragraph(3), Created: time.Now().UTC()}},
		{task: togo.Task{Name: f.Person().Name(), Description: f.Lorem().Paragraph(3), Created: time.Now().UTC(), Completed: daysFromNow(1)}},
		{task: togo.Task{Name: f.Person().Name(), Description: f.Lorem().Paragraph(3), Created: time.Now().UTC(), Completed: daysFromNow(2)}, dueDate: daysFromNow(3)},
	}

	for _, testCase := range testCases {
		expected := testCase.task
		if testCase.dueDate != nil {
			expected.AddDueDate(*testCase.dueDate)
		}

		err := pg.AddOrUpdateTask(expected)
		if err != nil {
			t.Error(err)
		}
		outcome, err := pg.FindTaskByName(expected.Name)

		if err != nil {
			t.Error("unable to find task by name")
		}

		err = compareTasks(expected, outcome)
		if err != nil {
			t.Error(err)
		}
	}
}

func compareTasks(expected, outcome togo.Task) error {

	if outcome.Name != expected.Name {
		return errors.New("names do not match")
	}

	if outcome.Description != expected.Description {
		return errors.New("descriptions do not match")
	}

	if outcome.Created.Equal(expected.Created) {
		return errors.New("creation dates do not match")
	}

	if outcome.DueOn() != expected.DueOn() {
		return errors.New("due dates do not match")
	}

	if outcome.Completed != expected.Completed {
		return errors.New("completion dates do not match")
	}

	return nil
}

func compareTime(expected *time.Time, outcome *time.Time) bool {
	if expected == nil && outcome == nil {
		return true
	}
	if expected == nil && outcome != nil {
		return false
	}
	if expected != nil && outcome == nil {
		return false
	}

	return expected.Equal(*outcome)
}
