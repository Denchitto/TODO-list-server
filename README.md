# Файлы для итогового задания

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.


фиксануть мейн
подруб постгреса
исп контекст и методы?
.env


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
			limit = 50
		)

		value := r.URL.Query().Get("search")

		args := make(map[string]any)
		if strings.TrimSpace(value) == "" {
			query = db.QueryByAllResults
			args["limit"] = limit
		} else if t, err := time.Parse(service.DateFormatSearch, value); err == nil {
			value = t.Format(service.DateFormat)
			query = db.QueryByDate
			args["limit"] = limit
			args["date"] = value
		} else {
			value = "%" + value + "%"
			query = db.QueryByTitleComment
			args["limit"] = limit
			args["search"] = value
		}

		tasks, err := db.GetTasks(srv.Db, query, args)
		if err != nil {
			writeJson(w, err, http.StatusInternalServerError, srv.Lg)
			return
		}
		writeJson(w, tasksResp{Tasks: tasks}, http.StatusOK, srv.Lg)
	}
}

func getTaskHandler(srv *server.CustomServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.URL.Query().Get("id")
		if strings.TrimSpace(id) == "" {
			writeJson(w, fmt.Errorf("id not specified"), http.StatusBadRequest, srv.Lg)
			return
		}
		if _, err := strconv.Atoi(id); err != nil {
			writeJson(w, fmt.Errorf("incorrect id format: %s; %w", id, err), http.StatusBadRequest, srv.Lg)
			return
		}

		task, err := db.GetTasks(srv.Db, db.QueryById, map[string]any{"id": id})
		if err != nil {
			writeJson(w, err, http.StatusInternalServerError, srv.Lg)
			return
		}
		if len(task) == 0 {
			writeJson(w, fmt.Errorf("unable to find task with id: %s", id), http.StatusNotFound, srv.Lg)
			return
		}
		writeJson(w, task[0], http.StatusOK, srv.Lg)
	}
}
func GetTasks(db *sql.DB, query string, args map[string]any) ([]*Task, error) {

	namedArgs := make([]any, 0, len(args))
	for k, v := range args {
		namedArgs = append(namedArgs, sql.Named(k, v))
	}

	rows, err := db.Query(query, namedArgs...)
	if err != nil {
		return []*Task{}, err
	}
	defer rows.Close()

	tasks := []*Task{}

	for rows.Next() {

		var newTask Task
		if err := rows.Scan(&newTask.ID,
			&newTask.Date,
			&newTask.Title,
			&newTask.Comment,
			&newTask.Repeat); err != nil {
			return []*Task{}, err
		}

		tasks = append(tasks, &newTask)
	}

	if err := rows.Err(); err != nil {
		return []*Task{}, err
	}
	return tasks, nil
}