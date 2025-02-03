package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"todo-app/pkg"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func main() {
	ctx := context.WithValue(context.Background(), "TraceID", pkg.GenerateTraceID())

	add := flag.String("add", "", "Add a new to-do item.")
	update := flag.String("update", "", "Update a to-do item, format ID:Description.")
	delete := flag.Int("delete", -1, "Delete a to-do item by ID.")
	list := flag.Bool("list", false, "List all to-do items.")
	status := flag.String("status", "", "Update the status of a to-do item, format ID:Status.")

	flag.Parse()

	todos, err := pkg.LoadTodos(ctx)
	if err != nil {
		logger.With("TraceID", ctx.Value("TraceID")).Error("Error loading to-dos", "error", err)
		return
	}

	switch {
	case *add != "":
		todos, err = pkg.AddTodo(ctx, todos, *add)
		if err != nil {
			logger.With("TraceID", ctx.Value("TraceID")).Error("Failed to add todo", "error", err)
			return
		}

	case *update != "":
		todos, err = pkg.UpdateTodoDescription(ctx, todos, *update)
		if err != nil {
			logger.With("TraceID", ctx.Value("TraceID")).Error("Failed to update description", "error", err)
			return
		}

	case *delete != -1:
		todos, err = pkg.DeleteTodo(ctx, todos, *delete)
		if err != nil {
			logger.With("TraceID", ctx.Value("TraceID")).Error("Failed to delete todo", "error", err)
			return
		}

	case *list:
		pkg.ListTodos(ctx, todos)
		return

	case *status != "":
		todos, err = pkg.UpdateTodoStatus(ctx, todos, *status)
		if err != nil {
			logger.With("TraceID", ctx.Value("TraceID")).Error("Failed to update status", "error", err)
			return
		}

	default:
		logger.With("TraceID", ctx.Value("TraceID")).Warn("No operation specified. Use --help for available commands.")
		return
	}
	err = pkg.SaveTodos(ctx, todos)
	if err != nil {
		logger.With("TraceID", ctx.Value("TraceID")).Error("Failed to save todos", "error", err)
	}
}
