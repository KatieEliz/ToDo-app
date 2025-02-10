package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"todo-app/pkg"

	"github.com/google/uuid"
)

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.New().String()
		ctx := context.WithValue(r.Context(), "TraceID", traceID)
		log.Printf("Request: %s %s [TraceID: %s]", r.Method, r.URL.Path, traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Get all todos
func GetTodosHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	todos, err := pkg.LoadTodos(ctx)
	if err != nil {
		http.Error(w, "Failed to load todos", http.StatusInternalServerError)
		log.Println("Error loading todos:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// Create a new todo
func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	description := r.PostFormValue("description")
	if description == "" {
		http.Error(w, "Description cannot be empty", http.StatusBadRequest)
		return
	}

	todos, err := pkg.LoadTodos(ctx)
	if err != nil {
		http.Error(w, "Failed to load todos", http.StatusInternalServerError)
		return
	}

	todos, err = pkg.AddTodo(ctx, todos, description)
	if err != nil {
		http.Error(w, "Failed to add todo", http.StatusInternalServerError)
		return
	}

	if err := pkg.SaveTodos(ctx, todos); err != nil {
		http.Error(w, "Failed to save todos", http.StatusInternalServerError)
		return
	}

	log.Printf("New to-do added: %s", description)
	http.Redirect(w, r, "http://localhost:8080/list", http.StatusSeeOther)
}

// Update a todo
func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input struct {
		ID          int    `json:"id"`
		Description string `json:"description,omitempty"`
		Status      string `json:"status,omitempty"`
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

	updated := false

	if input.Description != "" {
		updateInput := fmt.Sprintf("%d:%s", input.ID, input.Description)
		todos, err = pkg.UpdateTodoDescription(ctx, todos, updateInput)
		if err != nil {
			http.Error(w, "Failed to update todo description", http.StatusInternalServerError)
			return
		}
		updated = true
	}

	if input.Status != "" {
		updateInput := fmt.Sprintf("%d:%s", input.ID, input.Status)
		todos, err = pkg.UpdateTodoStatus(ctx, todos, updateInput)
		if err != nil {
			http.Error(w, "Failed to update todo status", http.StatusInternalServerError)
			return
		}
		updated = true
	}

	if !updated {
		http.Error(w, "No valid update provided", http.StatusBadRequest)
		return
	}

	if err := pkg.SaveTodos(ctx, todos); err != nil {
		http.Error(w, "Failed to save updated todos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// Delete a todo
func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input struct {
		ID int `json:"id"`
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

	todos, err = pkg.DeleteTodo(ctx, todos, input.ID)
	if err != nil {
		http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
		return
	}

	if err := pkg.SaveTodos(ctx, todos); err != nil {
		http.Error(w, "Failed to save todos after deletion", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// Start API Server
func StartServer() {
	mux := http.NewServeMux()
	mux.Handle("/create", TraceMiddleware(http.HandlerFunc(CreateTodoHandler)))
	mux.Handle("/get", TraceMiddleware(http.HandlerFunc(GetTodosHandler)))
	mux.Handle("/update", TraceMiddleware(http.HandlerFunc(UpdateTodoHandler)))
	mux.Handle("/delete", TraceMiddleware(http.HandlerFunc(DeleteTodoHandler)))

	log.Println("Starting API server on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
