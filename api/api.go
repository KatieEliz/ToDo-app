package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"todo-app/pkg"

	"github.com/google/uuid"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func traceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.New().String()
		ctx := context.WithValue(r.Context(), "TraceID", traceID)
		logger.With("TraceID", traceID).Info("Incoming request", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	todos, err := pkg.LoadTodos(ctx)
	if err != nil {
		http.Error(w, "Failed to load todos", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(todos)
}
func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input struct {
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	todos, err := pkg.LoadTodos(ctx)
	if err != nil {
		http.Error(w, "Failed to load todos", http.StatusInternalServerError)
		return
	}

	todos, err = pkg.AddTodo(ctx, todos, input.Description)
	if err != nil {
		http.Error(w, "Failed to add todo", http.StatusInternalServerError)
		return
	}

	_ = pkg.SaveTodos(ctx, todos)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todos)
}
