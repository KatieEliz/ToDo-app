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
	pkg.ListTodos(todos)
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

	todos = pkg.DeleteTodo(todos, 1)
	fmt.Println("\nTodos after deleting item with ID 1:")
	for _, todo := range todos {
		fmt.Printf("ID: %d, Description: %s, Status: %s\n", todo.ID, todo.Description, todo.Status)
	}

	err = pkg.SaveTodos(todos)
	if err != nil {
		log.Fatal("Error saving todos after deletion:", err)
	}
}
