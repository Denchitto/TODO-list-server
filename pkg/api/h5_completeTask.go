package api

import (
	"database/sql"
	"errors"
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

		task, err := db.ReadTask(srv.Db, db.QueryById, id)
		if errors.Is(sql.ErrNoRows, err) {
			writeJson(w, fmt.Errorf("unable to find task with id: %s", id), http.StatusNotFound, srv.Lg)
			return
		}
		if err != nil {
			writeJson(w, err, http.StatusInternalServerError, srv.Lg)
			return
		}

		if strings.TrimSpace(task.Repeat) == "" {
			if _, err := db.DeleteTask(srv.Db, db.QueryDeleteTask, task.ID); err != nil {
				writeJson(w, err, http.StatusInternalServerError, srv.Lg)
				return
			}

			srv.Lg.Printf("task %s succesfully deleted", task.Title)

			writeJson(w, struct{}{}, http.StatusOK, srv.Lg)
			return
		} else {
			task.Date, err = service.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				writeJson(w, err, http.StatusInternalServerError, srv.Lg)
				return
			}
			if _, err = db.UpdateTask(srv.Db, db.QueryUpdateTaskById, task.Date, task.ID); err != nil {
				writeJson(w, err, http.StatusInternalServerError, srv.Lg)
				return
			}

			srv.Lg.Printf("task %s succesfully updated", task.Title)

			writeJson(w, struct{}{}, http.StatusOK, srv.Lg)
		}
	}
}
