package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func Init(db *sql.DB, lg *log.Logger) (*http.ServeMux, error) {

	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting working directory: %w", err)
	}
	webDir := filepath.Join(wd, "web")

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	// Прокидываем логгер и бд в хендлеры
	mux.HandleFunc("/api/signin", authorization(lg))
	mux.HandleFunc("/api/nextdate", nextDateHandler(lg))
	mux.HandleFunc("/api/tasks", auth(getTasksHandler(db, lg), lg))
	mux.HandleFunc("/api/task/done", auth(completeTaskHandler(db, lg), lg))
	mux.HandleFunc("/api/task", auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			addTaskHandler(db, lg)(w, r)
		case http.MethodGet:
			getTaskHandler(db, lg)(w, r)
		case http.MethodPut:
			editTaskHandler(db, lg)(w, r)
		case http.MethodDelete:
			deleteTaskHandler(db, lg)(w, r)
		default:
			lg.Printf("method not allowed: %s", r.Method)
			writeJson(w, fmt.Errorf("method not allowed"), http.StatusMethodNotAllowed, lg)
		}
	}, lg))
	return mux, nil
}
