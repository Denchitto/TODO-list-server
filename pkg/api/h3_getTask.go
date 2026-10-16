package api

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Denchitto/TODO-list-server/pkg/db"
)

type tasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

const limit = 50

func getTasksHandler(database *sql.DB, lg *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		value := r.URL.Query().Get("search")

		tasks, err := db.ReadTasks(database, value, limit)
		if err != nil {
			lg.Println(err)
			writeJson(w, fmt.Errorf("internal error"), http.StatusInternalServerError, lg)
			return
		}

		lg.Println("successfully received list of tasks")

		writeJson(w, tasksResp{Tasks: tasks}, http.StatusOK, lg)
	}
}

func getTaskHandler(database *sql.DB, lg *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.URL.Query().Get("id")

		if strings.TrimSpace(id) == "" {
			writeJson(w, fmt.Errorf("id is not specified"), http.StatusBadRequest, lg)
			return
		}
		if _, err := strconv.Atoi(id); err != nil {
			writeJson(w, fmt.Errorf("incorrect id format: %s; %w", id, err), http.StatusBadRequest, lg)
			return
		}

		task, err := db.ReadTask(database, []any{id}...)
		if errors.Is(sql.ErrNoRows, err) {
			writeJson(w, fmt.Errorf("unable to find task with id: %s", id), http.StatusNotFound, lg)
			return
		}
		if err != nil {
			lg.Println(err)
			writeJson(w, fmt.Errorf("internal error"), http.StatusInternalServerError, lg)
			return
		}

		lg.Printf("task %s succesfully received", id)

		writeJson(w, task, http.StatusOK, lg)
	}
}
