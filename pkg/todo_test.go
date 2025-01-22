package pkg

import (
	"os"
	"testing"
)

func TestCleanUpOfExistingTodoJsonFiles(t *testing.T) {
	err := os.Remove("todo.json")
	if err != nil && !os.IsNotExist(err) {

		t.Fatalf("Failed to clean up todo.json: %v", err)
	}
}

func TestLoadingTodosWhenFileDoesntExist(t *testing.T) {
	_, err := LoadTodos()
	if err != nil {
		t.Fatalf("Failed to load todos: %v", err)
	}
}

func TestReturnedSliceIsEmptyWhenFileDoesntExist(t *testing.T) {
	todos, _ := LoadTodos()
	if len(todos) != 0 {
		t.Errorf("Expected 0 todos, got %d", len(todos))
	}
}

func TestCreateSampleTodoFileAndSaveSampleToFile(t *testing.T) {
	expectedTodos := []TodoItem{
		{ID: 1, Description: "Test todo", Status: "pending"},
	}

	err := SaveTodos(expectedTodos)
	if err != nil {
		t.Fatalf("Failed to save todos: %v", err)
	}
}
