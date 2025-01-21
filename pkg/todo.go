package make

import (
	"encoding/json"
	"os"
)

type Todo struct {
	Description string
}

func LoadTodos() ([]Todo, error) {
	file, err := os.Open("todos.json") //try to open file, error if doesnt exist

	if err != nil { //if error, checks if exists using isnotexist
		if os.IsNotExist(err) {
			return []Todo{}, nil //doesnt exist empty list of todos
		}
		return nil, err
	}
	defer file.Close() //function finish = file closed

	var todos []Todo
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&todos)

	if err != nil {
		return nil, err
	}

	return todos, nil

}
