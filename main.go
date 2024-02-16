package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const port = ":8000"

type PageVariables struct {
	Title string
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		name := request.URL.Query().Get("name")
		if name == "" {
			name = "No Name is defined."
		}
		response := fmt.Sprintf("<html><body><h1>Welcome, %s</h1></body></html>", name)
		writer.Header().Set("Content-Type", "text/html")
		writer.Write([]byte(response))
	})
	router.Get("/html", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "index", PageVariables{Title: "Hello,World!"})
	})
	fmt.Printf("Listening on port%s\n", port)
	http.ListenAndServe(port, router)
}

func renderTemplate(writer http.ResponseWriter, tmpl string, data PageVariables) {
	tmplFiles := []string{"templates/layout.html", "templates/" + tmpl + ".html"}
	t, err := template.ParseFiles(tmplFiles...)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(writer, "layout", data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
