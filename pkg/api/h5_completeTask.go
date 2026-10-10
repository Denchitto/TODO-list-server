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

func completeTaskHandler(srv *server.CustomServer) http.HandlerFunc {
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

		task, err := db.ReadTasks(srv.Db, db.QueryById, id)
		if err != nil {
			writeJson(w, err, http.StatusInternalServerError, srv.Lg)
			return
		}
		if len(task) == 0 {
			writeJson(w, fmt.Errorf("unable to find task with id: %s", id), http.StatusNotFound, srv.Lg)
			return
		}

		if strings.TrimSpace(task[0].Repeat) == "" {
			if _, err := db.DeleteTask(srv.Db, db.QueryDeleteTask, task[0].ID); err != nil {
				writeJson(w, err, http.StatusInternalServerError, srv.Lg)
				return
			}
			srv.Lg.Printf("task %s succesfully deleted", task[0].Title)
			writeJson(w, struct{}{}, http.StatusOK, srv.Lg)
			return
		} else {
			task[0].Date, err = service.NextDate(time.Now(), task[0].Date, task[0].Repeat)
			if err != nil {
				writeJson(w, err, http.StatusInternalServerError, srv.Lg)
				return
			}
			if _, err = db.UpdateTask(srv.Db, db.QueryUpdateTaskById,
				task[0].Date,
				task[0].ID); err != nil {
				writeJson(w, err, http.StatusInternalServerError, srv.Lg)
				return
			}
			srv.Lg.Printf("task %s succesfully updated", task[0].Title)
			writeJson(w, struct{}{}, http.StatusOK, srv.Lg)
		}
	}
}
