package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"todo-app/api"
	"todo-app/pkg"
)

var static = os.DirFS("static")

func StartServer() {
	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Combined Page (About and List)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Load todos for the list
		todos, err := pkg.LoadTodos(ctx)
		if err != nil {
			http.Error(w, "Failed to load todos", http.StatusInternalServerError)
			log.Println("Error loading todos:", err)
			return
		}

		// Parse both About and Todo List template
		tmpl, err := template.ParseFiles("static/combined.html") // Single template for both pages
		if err != nil {
			http.Error(w, "Unable to load template", http.StatusInternalServerError)
			log.Println("Error loading template:", err)
			return
		}

		// Execute template with todos and about information
		err = tmpl.Execute(w, struct {
			Todos []pkg.TodoItem
			About string
		}{
			Todos: todos, // Only pass the todos slice here
			About: "This is a simple to-do application built with Go. Manage your tasks easily and track their progress!",
		})
		if err != nil {
			http.Error(w, "Unable to render template", http.StatusInternalServerError)
			log.Println("Error executing template:", err)
		}
	})

	// Redirect API endpoints to API server
	mux.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://localhost:8081/create", http.StatusTemporaryRedirect)
	})
	mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://localhost:8081/get", http.StatusTemporaryRedirect)
	})
	mux.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://localhost:8081/update", http.StatusTemporaryRedirect)
	})
	mux.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://localhost:8081/delete", http.StatusTemporaryRedirect)
	})

	log.Println("Starting web server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func main() {
	isAPI := flag.Bool("api", false, "Run the API server")
	flag.Parse()

	if *isAPI {
		api.StartServer() // Run the API server on port 8081
		return
	}

	StartServer() // Run the web server on port 8080
}
