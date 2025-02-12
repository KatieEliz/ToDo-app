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
		ctx := r.Context()
		ctx = context.WithValue(ctx, "TraceID", traceID)
		log.Printf("Request: %s %s [TraceID: %s]", r.Method, r.URL.Path, traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetTodosHandler(w http.ResponseWriter, r *http.Request) {
	todoStore := pkg.NewTodoStore()
	todos, err := todoStore.LoadTodos()
	if err != nil {
		http.Error(w, "Failed to load todos", http.StatusInternalServerError)
		log.Println("Error loading todos:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	description := r.PostFormValue("description")
	if description == "" {
		http.Error(w, "Description cannot be empty", http.StatusBadRequest)
		return
	}

	todoStore := pkg.NewTodoStore()
	_, err := todoStore.LoadTodos()
	if err != nil {
		http.Error(w, "Failed to load todos", http.StatusInternalServerError)
		return
	}

	todoStore.AddTodo(description)

	if err := todoStore.SaveTodos(); err != nil {
		http.Error(w, "Failed to save todos", http.StatusInternalServerError)
		return
	}

	log.Printf("New to-do added: %s", description)
	http.Redirect(w, r, "http://localhost:8080/list", http.StatusSeeOther)
}

func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID          int    `json:"id"`
		Description string `json:"description,omitempty"`
		Status      string `json:"status,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	todoStore := pkg.NewTodoStore()
	todos, err := todoStore.LoadTodos()
	if err != nil {
		http.Error(w, "Failed to load todos", http.StatusInternalServerError)
		return
	}

	updated := false

	if input.Description != "" {
		updateInput := fmt.Sprintf("%d:%s", input.ID, input.Description)
		todoStore.UpdateTodoDescription(updateInput)
		updated = true
	}

	if input.Status != "" {
		updateInput := fmt.Sprintf("%d:%s", input.ID, input.Status)
		todoStore.UpdateTodoStatus(updateInput)
		updated = true
	}

	if !updated {
		http.Error(w, "No valid update provided", http.StatusBadRequest)
		return
	}

	if err := todoStore.SaveTodos(); err != nil {
		http.Error(w, "Failed to save updated todos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID int `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	todoStore := pkg.NewTodoStore()
	todos, err := todoStore.LoadTodos()
	if err != nil {
		http.Error(w, "Failed to load todos", http.StatusInternalServerError)
		return
	}

	todoStore.DeleteTodo(input.ID)

	if err := todoStore.SaveTodos(); err != nil {
		http.Error(w, "Failed to save todos after deletion", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func StartServer() {
	mux := http.NewServeMux()
	mux.Handle("/create", TraceMiddleware(http.HandlerFunc(CreateTodoHandler)))
	mux.Handle("/get", TraceMiddleware(http.HandlerFunc(GetTodosHandler)))
	mux.Handle("/update", TraceMiddleware(http.HandlerFunc(UpdateTodoHandler)))
	mux.Handle("/delete", TraceMiddleware(http.HandlerFunc(DeleteTodoHandler)))

	log.Println("Starting API server on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
