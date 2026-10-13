package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Denchitto/TODO-list-server/pkg/server"
)

func NewRouter() (*http.ServeMux, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting working directory: %w", err)
	}
	webDir := filepath.Join(wd, "web")

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	return mux, nil
}

func Init(mux *http.ServeMux, srv *server.CustomServer) {
	mux.HandleFunc("/api/signin", authorization(srv.Lg))
	mux.HandleFunc("/api/nextdate", nextDateHandler(srv.Lg))
	mux.HandleFunc("/api/tasks", auth(getTasksHandler(srv), srv.Lg))
	mux.HandleFunc("/api/task/done", auth(completeTaskHandler(srv), srv.Lg))
	mux.HandleFunc("/api/task", auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			addTaskHandler(srv)(w, r)
		case http.MethodGet:
			getTaskHandler(srv)(w, r)
		case http.MethodPut:
			editTaskHandler(srv)(w, r)
		case http.MethodDelete:
			deleteTaskHandler(srv)(w, r)
		default:
			srv.Lg.Printf("method not supported: %s", r.Method)
		}
	}, srv.Lg))
}
