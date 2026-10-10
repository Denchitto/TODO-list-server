package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Denchitto/TODO-list-server/pkg/db"
	"github.com/Denchitto/TODO-list-server/pkg/server"
	"github.com/Denchitto/TODO-list-server/pkg/service"
)

type tasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func getTasksHandler(srv *server.CustomServer) http.HandlerFunc {
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

		tasks, err := db.ReadTasks(srv.Db, query, args...)
		if err != nil {
			writeJson(w, err, http.StatusInternalServerError, srv.Lg)
			return
		}

		srv.Lg.Println("succefully received list of tasks")

		writeJson(w, tasksResp{Tasks: tasks}, http.StatusOK, srv.Lg)
	}
}

func getTaskHandler(srv *server.CustomServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.URL.Query().Get("id")

		if strings.TrimSpace(id) == "" {
			writeJson(w, fmt.Errorf("id is not specified"), http.StatusBadRequest, srv.Lg)
			return
		}
		if _, err := strconv.Atoi(id); err != nil {
			writeJson(w, fmt.Errorf("incorrect id format: %s; %w", id, err), http.StatusBadRequest, srv.Lg)
			return
		}

		task, err := db.ReadTasks(srv.Db, db.QueryById, []any{id}...)
		if err != nil {
			writeJson(w, err, http.StatusInternalServerError, srv.Lg)
			return
		}
		if len(task) == 0 {
			writeJson(w, fmt.Errorf("unable to find task with id: %s", id), http.StatusNotFound, srv.Lg)
			return
		}

		srv.Lg.Printf("task %s succesfully received", id)

		writeJson(w, task[0], http.StatusOK, srv.Lg)
	}
}
