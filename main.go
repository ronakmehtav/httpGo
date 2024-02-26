package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

const port = ":8000"

type TodoItem struct {
	id     int
	Label  string
	Status bool
}

type PageVariables struct {
	Title     string
	TodoItems []TodoItem
}

const (
	FALSE = iota // 0
	TRUE
)

func main() {
	defaultItems := []TodoItem{}

	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlStmt := `
    CREATE TABLE IF NOT EXISTS task (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,label TEXT, status INTEGER);
    `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		fmt.Println("error with stmt")
		log.Fatal(err)
	}

	rows, err := db.Query("select * from task;")
	if err != nil {
		fmt.Println("error with the fetch")
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var task string
		var status int
		err := rows.Scan(&id, &task, &status)
		if err != nil {
			log.Fatal(err)
		}
		defaultItems = append(defaultItems, TodoItem{id: id, Label: task, Status: status == TRUE})
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/styles/{stylesPath}", func(w http.ResponseWriter, r *http.Request) {
		stylePath := chi.URLParam(r, "stylesPath")
		http.ServeFile(w, r, filepath.Join("./styles/", stylePath))
	})
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "index", PageVariables{Title: "ToDo App", TodoItems: defaultItems})
	})
	router.Get("/home", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "index", PageVariables{Title: "ToDo App", TodoItems: defaultItems})
	})

	router.Post("/api/addTask", func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength == 0 {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprintf(w, "no data received")
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse the form data", http.StatusBadRequest)
			return
		}
		task := strings.TrimSpace(r.FormValue("addTodo"))

		if len(task) == 0 {
			http.Redirect(w, r, "/", http.StatusNoContent)
			return
		}

		stmtquery := fmt.Sprintf("INSERT INTO task (label,status) values(\"%s\",%d);", task, FALSE)
		// fmt.Println(stmtquery)
		_, err := db.Exec(stmtquery)
		if err != nil {
			http.Error(w, "Failed to enter data into db.", http.StatusBadRequest)
			return
		}

		defaultItems = append(defaultItems, TodoItem{Label: task, Status: false})
		http.Redirect(w, r, "/", http.StatusSeeOther)

	})

	router.Put("/api/update/{index}", func(w http.ResponseWriter, r *http.Request) {
		index, err := strconv.Atoi(chi.URLParam(r, "index"))
		if err != nil {
			http.Error(w, "The value passed must be a number.", http.StatusBadRequest)
			return
		}

		if index < 0 || index >= len(defaultItems) {
			http.Error(w, "Incorrect value Passed.", http.StatusBadRequest)
			return
		}

		var stmtquery string
		if defaultItems[index].Status {
			stmtquery = fmt.Sprintf("UPDATE task SET status=%d WHERE id=%d;", FALSE, defaultItems[index].id)
		} else {
			stmtquery = fmt.Sprintf("UPDATE task SET status=%d WHERE id=%d;", TRUE, defaultItems[index].id)
		}
		_, err = db.Exec(stmtquery)
		if err != nil {
			http.Error(w, "Failed to enter data into db.", http.StatusBadRequest)
			return
		}
		defaultItems[index].Status = !defaultItems[index].Status
		w.WriteHeader(http.StatusOK)
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
