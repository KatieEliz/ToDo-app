package main

import (
	"fmt"
	"log"
	"todo-app/pkg"
)

func main() {
	todos, err := pkg.LoadTodos()
	if err != nil {
		log.Fatal("Error loading todos:", err)
	}

	fmt.Println("Loaded Todos:")
	for _, todo := range todos {
		fmt.Printf("ID: %d, Description: %s, Status: %s\n", todo.ID, todo.Description, todo.Status)
	}

	todos = pkg.AddTodo(todos, "Test new task")
	err = pkg.SaveTodos(todos)
	if err != nil {
		log.Fatal("Error saving todos:", err)
	}

	fmt.Println("\nTodos after adding a new item:")
	for _, todo := range todos {
		fmt.Printf("ID: %d, Description: %s, Status: %s\n", todo.ID, todo.Description, todo.Status)
	}

	pkg.UpdateTodoDescription(todos, "1:Updated task description")

}
