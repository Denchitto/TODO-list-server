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
	mux.HandleFunc("/api/nextdate", nextDayHandler)
	mux.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			addTaskHandler(srv)(w, r)
		}
	})
}
