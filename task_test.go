package togo

import "testing"

func TestNewTaskIsNotCompleted(t *testing.T) {
	task := NewTask("name", "description")

	if task.Completed() {
		t.Error("a new task should not be completed")
	}
}
