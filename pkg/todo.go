package pkg

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
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

func GenerateTraceID() string {
	return uuid.New().String()
}

func LoadTodos(ctx context.Context) ([]TodoItem, error) {
	traceIDValue := ctx.Value("TraceID")
	traceID, ok := traceIDValue.(string)
	if !ok || traceID == "" {
		traceID = "unknown"
	}

	// Open the todo.json file
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			logger.With("TraceID", traceID).Info("No existing todo file found. Starting fresh.")
			return []TodoItem{}, nil
		}
		logger.With("TraceID", traceID).Error("Failed to open todo file", "error", err)
		return nil, err
	}
	defer file.Close()

	// Read the entire file content
	data, err := io.ReadAll(file)
	if err != nil {
		logger.With("TraceID", traceID).Error("Failed to read todo file", "error", err)
		return nil, err
	}

	// Parse the JSON data into a slice of TodoItem
	var todos []TodoItem
	if err := json.Unmarshal(data, &todos); err != nil {
		logger.With("TraceID", traceID).Error("Failed to decode todo file", "error", err)
		return nil, err
	}

	logger.With("TraceID", traceID).Info("Loaded todos from file", "count", len(todos))
	return todos, nil
}

func SaveTodos(ctx context.Context, todos []TodoItem) error {
	traceID := ctx.Value("TraceID").(string)
	file, err := os.Create(filename)
	if err != nil {
		logger.With("TraceID", traceID).Error("Failed to open file for writing", "error", err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(todos)
	if err != nil {
		logger.With("TraceID", traceID).Error("Failed to encode JSON", "error", err)
		return err
	}

	logger.With("TraceID", traceID).Info("Saved todos successfully", "count", len(todos))
	return nil
}

func AddTodo(ctx context.Context, todos []TodoItem, description string) ([]TodoItem, error) {
	traceID := ctx.Value("TraceID").(string)

	id := 1
	for _, todo := range todos {
		if todo.ID >= id {
			id = todo.ID + 1
		}
	}

	newTodo := TodoItem{
		ID:          id,
		Description: description,
		Status:      "not started",
	}

	todos = append(todos, newTodo)
	logger.With("TraceID", traceID).Info("To-do item added", "id", id, "description", description)

	return todos, nil
}

func UpdateTodoDescription(ctx context.Context, todos []TodoItem, input string) ([]TodoItem, error) {
	traceID := ctx.Value("TraceID").(string)

	parts := strings.SplitN(input, ":", 2)
	if len(parts) != 2 {
		logger.With("TraceID", traceID).Warn("Invalid format for update", "input", input)
		return todos, nil
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil {
		logger.With("TraceID", traceID).Error("Invalid ID for update", "input", input, "error", err)
		return todos, err
	}

	description := parts[1]

	for index, todo := range todos {
		if todo.ID == id {
			todos[index].Description = description
			logger.With("TraceID", traceID).Info("To-do description updated", "id", id, "new_description", description)
			return todos, nil
		}
	}
	logger.With("TraceID", traceID).Warn("To-do item not found for update", "id", id)
	return todos, nil
}

func UpdateTodoStatus(ctx context.Context, todos []TodoItem, input string) ([]TodoItem, error) {
	traceID := ctx.Value("TraceID").(string)

	parts := strings.SplitN(input, ":", 2)
	if len(parts) != 2 {
		logger.With("TraceID", traceID).Warn("Invalid format for status update", "input", input)
		return todos, nil
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil {
		logger.With("TraceID", traceID).Error("Invalid ID for status update", "input", input, "error", err)
		return todos, err
	}

	newStatus := parts[1]

	for index, todo := range todos {
		if todo.ID == id {
			todos[index].Status = newStatus
			logger.With("TraceID", traceID).Info("To-do status updated", "id", id, "new_status", newStatus)
			return todos, nil
		}
	}
	logger.With("TraceID", traceID).Warn("To-do item not found for status update", "id", id)
	return todos, nil
}

func DeleteTodo(ctx context.Context, todos []TodoItem, id int) ([]TodoItem, error) {
	traceID := ctx.Value("TraceID").(string)

	for index, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:index], todos[index+1:]...)
			logger.With("TraceID", traceID).Info("To-do item deleted", "id", id)
			return todos, nil
		}
	}

	logger.With("TraceID", traceID).Warn("To-do item not found for deletion", "id", id)
	return todos, nil
}

func ListTodos(ctx context.Context, todos []TodoItem) {
	traceID := ctx.Value("TraceID").(string)

	if len(todos) == 0 {
		logger.With("TraceID", traceID).Info("No to-do items found.")
		return
	}

	logger.With("TraceID", traceID).Info("Listing to-do items", "count", len(todos))
	for _, todo := range todos {
		logger.With("TraceID", traceID).Info("To-do item", "id", todo.ID, "description", todo.Description, "status", todo.Status)
	}
}
