package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"todo-app/api"
	"todo-app/pkg"
)

var static = os.DirFS("static")

func StartServer() {
	todoStore := pkg.NewTodoStore()
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		todos, err := todoStore.LoadTodos()
		if err != nil {
			http.Error(w, "Failed to load todos", http.StatusInternalServerError)
			log.Println("Error loading todos:", err)
			return
		}

		tmpl, err := template.ParseFiles("static/combined.html")
		if err != nil {
			http.Error(w, "Unable to load template", http.StatusInternalServerError)
			log.Println("Error loading template:", err)
			return
		}

		err = tmpl.Execute(w, struct {
			Todos []pkg.TodoItem
			About string
		}{
			Todos: todos,
			About: "This is a simple to-do application built with Go. Manage your tasks easily and track their progress!",
		})
		if err != nil {
			http.Error(w, "Unable to render template", http.StatusInternalServerError)
			log.Println("Error executing template:", err)
		}
	})

	mux.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			description := r.PostFormValue("description")
			if description == "" {
				http.Error(w, "Description is required", http.StatusBadRequest)
				return
			}

			todoStore.AddTodo(description)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	mux.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			idStr := r.PostFormValue("id")
			description := r.PostFormValue("description")

			if idStr == "" || description == "" {
				http.Error(w, "ID and description are required", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid ID format", http.StatusBadRequest)
				return
			}

			todoStore.UpdateTodoDescription(fmt.Sprintf("%d:%s", id, description))
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	mux.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			idStr := r.PostFormValue("id")
			if idStr == "" {
				http.Error(w, "ID is required", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid ID format", http.StatusBadRequest)
				return
			}

			todoStore.DeleteTodo(id)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	log.Println("Starting web server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func main() {
	isAPI := flag.Bool("api", false, "Run the API server")
	flag.Parse()

	if *isAPI {
		api.StartServer()
		return
	}

	StartServer()
}
