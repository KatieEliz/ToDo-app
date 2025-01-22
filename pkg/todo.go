package main

import (
	"encoding/json"
	"os"
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
