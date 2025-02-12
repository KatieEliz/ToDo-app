package pkg

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"
)

type TodoItem struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

var (
	filename = "todo.json"
	logger   = slog.New(slog.NewTextHandler(os.Stdout, nil))
)

type TodoStore struct {
	mu    sync.Mutex
	todos []TodoItem
	ch    chan func()
}

func NewTodoStore() *TodoStore {
	ts := &TodoStore{
		todos: []TodoItem{},
		ch:    make(chan func()),
	}
	go ts.run()
	return ts
}

func (ts *TodoStore) run() {
	for op := range ts.ch {
		op()
	}
}

func (ts *TodoStore) LoadTodos() ([]TodoItem, error) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Info("No existing todo file found. Starting fresh.")
			return []TodoItem{}, nil
		}
		logger.Error("Failed to open todo file", "error", err)
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		logger.Error("Failed to read todo file", "error", err)
		return nil, err
	}

	var todos []TodoItem
	if err := json.Unmarshal(data, &todos); err != nil {
		logger.Error("Failed to decode todo file", "error", err)
		return nil, err
	}

	return todos, nil
}

func (ts *TodoStore) SaveTodos() error {
	ts.ch <- func() {
		file, err := os.Create(filename)
		if err != nil {
			logger.Error("Failed to open file for writing", "error", err)
			return
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")

		if err := encoder.Encode(ts.todos); err != nil {
			logger.Error("Failed to encode JSON", "error", err)
			return
		}

		logger.Info("Saved todos successfully", "count", len(ts.todos))
	}
	return nil
}

func (ts *TodoStore) AddTodo(description string) {
	ts.ch <- func() {
		id := 1
		for _, todo := range ts.todos {
			if todo.ID >= id {
				id = todo.ID + 1
			}
		}

		ts.todos = append(ts.todos, TodoItem{
			ID:          id,
			Description: description,
			Status:      "not started",
		})

		logger.Info("To-do item added", "id", id, "description", description)
		ts.SaveTodos()
	}
}

func (ts *TodoStore) UpdateTodoDescription(input string) {
	ts.ch <- func() {
		parts := strings.SplitN(input, ":", 2)
		if len(parts) != 2 {
			logger.Warn("Invalid format for update", "input", input)
			return
		}
		id, err := strconv.Atoi(parts[0])
		if err != nil {
			logger.Error("Invalid ID for update", "input", input, "error", err)
			return
		}
		description := parts[1]

		for index, todo := range ts.todos {
			if todo.ID == id {
				ts.todos[index].Description = description
				logger.Info("To-do description updated", "id", id, "new_description", description)
				ts.SaveTodos()
				return
			}
		}
		logger.Warn("To-do item not found for update", "id", id)
	}
}

func (ts *TodoStore) UpdateTodoStatus(input string) {
	ts.ch <- func() {
		parts := strings.SplitN(input, ":", 2)
		if len(parts) != 2 {
			logger.Warn("Invalid format for update", "input", input)
			return
		}
		id, err := strconv.Atoi(parts[0])
		if err != nil {
			logger.Error("Invalid ID for update", "input", input, "error", err)
			return
		}
		status := parts[1]

		for index, todo := range ts.todos {
			if todo.ID == id {
				ts.todos[index].Status = status
				logger.Info("To-do status updated", "id", id, "new_status", status)
				ts.SaveTodos()
				return
			}
		}
		logger.Warn("To-do item not found for update", "id", id)
	}
}

func (ts *TodoStore) DeleteTodo(id int) {
	ts.ch <- func() {
		for index, todo := range ts.todos {
			if todo.ID == id {
				ts.todos = append(ts.todos[:index], ts.todos[index+1:]...)
				logger.Info("To-do item deleted", "id", id)
				ts.SaveTodos()
				return
			}
		}
		logger.Warn("To-do item not found for deletion", "id", id)
	}
}
