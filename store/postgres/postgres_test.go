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
	theTime := time.Now().Add(duration)
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

	task := togo.NewTask(taskName, f.Lorem().Paragraph(3))
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
	t.Cleanup(func() {
		err := pg.RemoveTask(taskName)
		if err != nil {
			return
		}
	})

	expected := togo.NewTask(taskName, f.Lorem().Paragraph(3))
	err := pg.AddOrUpdateTask(expected)
	if err != nil {
		t.Error(err)
	}

	outcome, err := pg.FindTaskByName(taskName)

	if err != nil {
		t.Error("unable to find task by name")
	}

	err = compareTasks(expected, outcome)
	if err != nil {
		t.Error(err)
	}
}

func TestTasksCanBeRetrievedByName(t *testing.T) {
	pg := NewPgStore(connectionString)
	f := faker.New()

	testCases := []struct {
		task    togo.Task
		dueDate *time.Time
	}{
		{task: togo.Task{Name: f.Person().Name(), Description: f.Lorem().Paragraph(3), Created: time.Now()}},
		{task: togo.Task{Name: f.Person().Name(), Description: f.Lorem().Paragraph(3), Created: time.Now(), Completed: daysFromNow(1)}},
		{task: togo.Task{Name: f.Person().Name(), Description: f.Lorem().Paragraph(3), Created: time.Now(), Completed: daysFromNow(2)}, dueDate: daysFromNow(3)},
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
		t.Cleanup(func() {
			err := pg.RemoveTask(testCase.task.Name)
			if err != nil {
				return
			}
		})
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

func TestTasksCanBeRetrievedByDueDate(t *testing.T) {
	pg := NewPgStore(connectionString)
	f := faker.New()

	var created = time.Now()

	tasks := []togo.Task{
		{Name: f.Person().Name(), Description: f.Lorem().Paragraph(3), Created: created, DueDate: &created},
		{Name: f.Person().Name(), Description: f.Lorem().Paragraph(3), Created: created, DueDate: &created},
		{Name: f.Person().Name(), Description: f.Lorem().Paragraph(3), Created: created, DueDate: &created},
		{Name: f.Person().Name(), Description: f.Lorem().Paragraph(3), Created: *daysFromNow(-1), Completed: daysFromNow(1)},
		{Name: f.Person().Name(), Description: f.Lorem().Paragraph(3), Created: *daysFromNow(-2), Completed: daysFromNow(2)},
	}

	for _, task := range tasks {
		err := pg.AddOrUpdateTask(task)
		t.Cleanup(func() {
			err := pg.RemoveTask(task.Name)
			if err != nil {
				return
			}
		})
		if err != nil {
			t.Error(err)
			return
		}
	}

	found, err := pg.FindTasksByDueDate(created)
	if err != nil {
		t.Error(err)
		return
	}

	if len(found) != 3 {
		t.Errorf("expected %d found %d", 3, len(found))
	}
}

func TestOverdueTasksCanBeRetrieved(t *testing.T) {
	pg := NewPgStore(connectionString)
	f := faker.New()

	var created = time.Now()
	var due = created.Add(-24 * time.Hour)

	tasks := []togo.Task{
		{Name: f.Person().Name(), Description: f.Lorem().Paragraph(3), Created: created, DueDate: &due},
		{Name: f.Person().Name(), Description: f.Lorem().Paragraph(3), Created: created, DueDate: &due},
		{Name: f.Person().Name(), Description: f.Lorem().Paragraph(3), Created: created, DueDate: &due},
		{Name: f.Person().Name(), Description: f.Lorem().Paragraph(3), Created: *daysFromNow(-1), Completed: daysFromNow(1)},
		{Name: f.Person().Name(), Description: f.Lorem().Paragraph(3), Created: *daysFromNow(-2), Completed: daysFromNow(2)},
	}

	for _, task := range tasks {
		err := pg.AddOrUpdateTask(task)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			_ = pg.RemoveTask(task.Name)
		})

	}

	found, err := pg.FindOverdueTasks()
	if err != nil {
		t.Error(err)
	}

	if len(found) != 3 {
		t.Errorf("expected %d found %d", 3, len(found))
	}
}

func compareTasks(expected, outcome togo.Task) error {

	if outcome.Name != expected.Name {
		return errors.New("names do not match")
	}

	if outcome.Description != expected.Description {
		return errors.New("descriptions do not match")
	}

	if !compareTime(&expected.Created, &outcome.Created) {
		return errors.New("creation dates do not match")
	}

	if !compareTime(expected.DueOn(), outcome.DueOn()) {
		return errors.New("due dates do not match")
	}

	if !compareTime(expected.Completed, outcome.Completed) {
		return errors.New("completion times do not match")
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

	return expected.UnixMilli() == outcome.UnixMilli()
}
