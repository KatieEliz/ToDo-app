package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type TodoItem struct {
	ID          int    `json: "id"`
	Description string `json:"description"`
	Status      string `json: "status"`
}

var (
	filename = "todo.json"
)

func LoadTodos() ([]TodoItem, error) {
	file, err := os.Open(filename)

	if err != nil {
		if os.IsNotExist(err) {
			return []TodoItem{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var todos []TodoItem
	err = json.NewDecoder(file).Decode(&todos)
	if err != nil {
		return nil, err
	}

	return todos, nil

}

func saveTodos(todos []TodoItem) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(todos)
}

func addTodo(todos []TodoItem, description string) []TodoItem {
	id := 1
	if len(todos) > 0 {
		id = todos[len(todos)-1].ID + 1
	}
	todos = append(todos, TodoItem{ID: id, Description: description, Status: "pending"})
	fmt.Println("To-do item added")
	return todos
}

func updateTodoDescription(todos []TodoItem, input string) {
	parts := strings.SplitN(input, ":", 2)
	if len(parts) != 2 {
		fmt.Println("Invalid format. Use ID:Description.")
		return
	}
}
