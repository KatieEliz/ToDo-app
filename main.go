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
	// Initialize the TodoStore
	todoStore := pkg.NewTodoStore()

	// Mux for routing
	mux := http.NewServeMux()

	// Serve static files
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
		api.StartServer()
		return
	}

	StartServer()
}
