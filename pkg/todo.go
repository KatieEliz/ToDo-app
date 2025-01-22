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

func LoadTodos() ([]TodoItem, error) {
	file, err := os.Open("todos.json") //try to open file, error if doesnt exist

	if err != nil { //if error, checks if exists using isnotexist
		if os.IsNotExist(err) {
			return []TodoItem{}, nil //doesnt exist empty list of todos
		}
		return nil, err
	}
	defer file.Close() //function finish = file close

	var todos []TodoItem
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&todos)

	if err != nil {
		return nil, err
	}

	return todos, nil

}
