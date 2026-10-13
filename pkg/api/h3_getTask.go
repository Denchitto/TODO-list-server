package api

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Denchitto/TODO-list-server/pkg/db"
	"github.com/Denchitto/TODO-list-server/pkg/service"
)

type tasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func getTasksHandler(database *sql.DB, lg *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var (
			query string
			args  []any
			limit = 50
		)

		value := r.URL.Query().Get("search")

		if strings.TrimSpace(value) == "" {
			query = db.QueryByAllResults
			args = []any{limit}
		} else if t, err := time.Parse(service.DateFormatSearch, value); err == nil {
			value = t.Format(service.DateFormat)
			query = db.QueryByDate
			args = []any{value, limit}
		} else {
			value = "%" + value + "%"
			query = db.QueryByTitleComment
			args = []any{value, value, limit}
		}

		tasks, err := db.ReadTasks(database, query, args...)
		if err != nil {
			writeJson(w, err, http.StatusInternalServerError, lg)
			return
		}

		lg.Println("succefully received list of tasks")

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

		task, err := db.ReadTask(database, db.QueryById, []any{id}...)
		if errors.Is(sql.ErrNoRows, err) {
			writeJson(w, fmt.Errorf("unable to find task with id: %s", id), http.StatusNotFound, lg)
			return
		}
		if err != nil {
			writeJson(w, err, http.StatusInternalServerError, lg)
			return
		}

		lg.Printf("task %s succesfully received", id)

		writeJson(w, task, http.StatusOK, lg)
	}
}
