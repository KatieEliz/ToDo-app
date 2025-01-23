package main

import (
	"flag"
	"fmt"
	"todo-app/pkg"
)

func main() {
	add := flag.String("add", "", "Add a new to-do item.")
	update := flag.String("update", "", "Update a to-do item, format ID:Description.")
	delete := flag.Int("delete", -1, "Delete a to-do item by ID.")
	list := flag.Bool("list", false, "List all to-do items.")
	status := flag.String("status", "", "Update the status of a to-do item, format ID:Status.")

	flag.Parse()

	todos, err := pkg.LoadTodos()
	if err != nil {
		fmt.Printf("Error loading to-dos: %v\n", err)
		return
	}

	switch {
	case *add != "":
		todos = pkg.AddTodo(todos, *add)
		pkg.SaveTodos(todos)

	case *update != "":
		pkg.UpdateTodoDescription(todos, *update)
		pkg.SaveTodos(todos)

	case *delete != -1:
		todos = pkg.DeleteTodo(todos, *delete)
		pkg.SaveTodos(todos)

	case *list:
		pkg.ListTodos(todos)

	case *status != "":
		pkg.UpdateTodoStatus(todos, *status)
		pkg.SaveTodos(todos)

	default:
		fmt.Println("No operation specified. Use --help for available commands.")
	}
}
