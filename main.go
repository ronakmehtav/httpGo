package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const port = ":8000"

// We would need page level. For the inital page creation.
// of all the value.
// should be of different type.

type TodoItem struct {
    Label string
    Status bool
}

type PageVariables struct {
	Title string
    TodoItems []TodoItem
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/styles/{stylesPath}", func(w http.ResponseWriter, r *http.Request) {
		stylePath := chi.URLParam(r, "stylesPath")
		http.ServeFile(w, r, filepath.Join("./styles/", stylePath))
	})
    defaultItems := []TodoItem{
        {Label: "hello, Word", Status: false},
        {Label: "study" ,Status:true},
        {Label: "eat" , Status: false},
    }
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "index", PageVariables{Title: "ToDo App", TodoItems: defaultItems })
	})
	router.Get("/home", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "index", PageVariables{Title: "ToDo App", TodoItems: defaultItems})
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
