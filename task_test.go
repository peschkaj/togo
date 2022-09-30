package togo

import "testing"

func TestNewTaskIsNotCompleted(t *testing.T) {
	task := NewTask("name", "description")

	if task.IsCompleted() {
		t.Error("a new task should not be Completed")
	}
}
