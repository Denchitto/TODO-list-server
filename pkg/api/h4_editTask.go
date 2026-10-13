package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Denchitto/TODO-list-server/pkg/db"
)

func editTaskHandler(database *sql.DB, lg *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var (
			task db.Task
			buf  bytes.Buffer
		)

		if _, err := buf.ReadFrom(r.Body); err != nil {
			writeJson(w, err, http.StatusBadRequest, lg)
			return
		}

		if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
			writeJson(w, err, http.StatusBadRequest, lg)
			return
		}

		if strings.TrimSpace(task.Title) == "" {
			writeJson(w, fmt.Errorf("empty title"), http.StatusBadRequest, lg)
			return
		}

		if err := checkDate(&task); err != nil {
			writeJson(w, err, http.StatusBadRequest, lg)
			return
		}

		count, err := db.UpdateTask(database, db.QueryUpdateTask, task.Date,
			task.Title,
			task.Comment,
			task.Repeat,
			task.ID)
		if err != nil {
			writeJson(w, err, http.StatusInternalServerError, lg)
			return
		} else if count == 0 {
			writeJson(w, fmt.Errorf("incorrect id for updating task"), http.StatusNotFound, lg)
			return
		}

		lg.Printf("task %s succesfully updated", task.Title)

		writeJson(w, struct{}{}, http.StatusOK, lg)
	}
}
