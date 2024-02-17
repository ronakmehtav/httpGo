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
	Name  string
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		userName := r.URL.Query().Get("name")
		renderTemplate(w, "index", PageVariables{Title: "Hello,World!", Name: userName})
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
