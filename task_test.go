package togo

import (
	"sort"
	"testing"
)

func TestNewTaskIsNotCompleted(t *testing.T) {
	task := NewTask("name", "description")

	if task.IsCompleted() {
		t.Error("a new task should not be Completed")
	}
}

func TestTasksSort(t *testing.T) {
	ts := []Task{}

	zero := Task{Name: "0", Priority: None}
	one := Task{Name: "1", Priority: None}
	two := Task{Name: "2", Priority: Low}
	three := Task{Name: "3", Priority: Low}
	four := Task{Name: "4", Priority: Medium}
	five := Task{Name: "5", Priority: High}

	ts = append(ts, two)
	ts = append(ts, zero)
	ts = append(ts, five)
	ts = append(ts, four)
	ts = append(ts, one)
	ts = append(ts, three)

	sort.Sort(Tasks(ts))

	if zero != ts[0] {
		t.Error("did not sort")
	}

	if five != ts[5] {
		t.Error("did not sort")
	}
}
