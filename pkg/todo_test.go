package pkg

import (
	"context"
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
	ctx := context.Background()
	_, err := LoadTodos(ctx)
	if err != nil {
		t.Fatalf("Failed to load todos: %v", err)
	}
}

func TestReturnedSliceIsEmptyWhenFileDoesntExist(t *testing.T) {
	ctx := context.Background()
	todos, err := LoadTodos(ctx)
	if err != nil {
		t.Fatalf("Unexpected error while loading todos: %v", err)
	}
	if len(todos) != 0 {
		t.Errorf("Expected 0 todos, got %d", len(todos))
	}
}

func TestCreateSampleTodoFileAndSaveSampleToFile(t *testing.T) {
	ctx := context.Background()
	expectedTodos := []TodoItem{
		{ID: 1, Description: "Test todo", Status: "pending"},
	}

	err := SaveTodos(ctx, expectedTodos)
	if err != nil {
		t.Fatalf("Failed to save todos: %v", err)
	}

	loadedTodos, err := LoadTodos(ctx)
	if err != nil {
		t.Fatalf("Failed to load saved todos: %v", err)
	}

	if len(loadedTodos) != len(expectedTodos) {
		t.Errorf("Expected %d todos, got %d", len(expectedTodos), len(loadedTodos))
	}

	_ = os.Remove("todo.json")
}
