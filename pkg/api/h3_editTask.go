package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Denchitto/TODO-list-server/pkg/db"
	"github.com/Denchitto/TODO-list-server/pkg/server"
)

func editTaskHandler(srv *server.CustomServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var (
			task db.Task
			buf  bytes.Buffer
		)

		if _, err := buf.ReadFrom(r.Body); err != nil {
			writeJson(w, err, http.StatusBadRequest, srv.Lg)
			return
		}

		if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
			writeJson(w, err, http.StatusBadRequest, srv.Lg)
			return
		}

		if strings.TrimSpace(task.Title) == "" {
			writeJson(w, fmt.Errorf("empty title"), http.StatusBadRequest, srv.Lg)
			return
		}

		if err := checkDate(&task); err != nil {
			writeJson(w, err, http.StatusBadRequest, srv.Lg)
			return
		}

		notFound, err := db.UpdateTask(srv.Db, db.QueryUpdateTask, task.Date,
			task.Title,
			task.Comment,
			task.Repeat,
			task.ID)
		if err != nil {
			writeJson(w, err, http.StatusInternalServerError, srv.Lg)
			return
		} else if notFound {
			writeJson(w, fmt.Errorf("incorrect id for updating task"), http.StatusNotFound, srv.Lg)
			return
		}

		srv.Lg.Printf("task %s succesfully updated", task.Title)

		writeJson(w, struct{}{}, http.StatusOK, srv.Lg)
	}
}
