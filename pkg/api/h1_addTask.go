package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Denchitto/TODO-list-server/pkg/db"
	"github.com/Denchitto/TODO-list-server/pkg/server"
	"github.com/Denchitto/TODO-list-server/pkg/service"
)

type miniID struct {
	ID string `json:"id"`
}

func checkDate(task *db.Task) error {

	now := time.Now()

	switch {

	case strings.TrimSpace(task.Date) == "":
		task.Date = now.Format(service.DateFormat)

	default:
		parsedDate, err := time.Parse(service.DateFormat, task.Date)
		if err != nil {
			return err
		}

		if service.AfterNow(now, parsedDate) {

			switch {

			case strings.TrimSpace(task.Repeat) == "":
				task.Date = now.Format(service.DateFormat)

			default:
				task.Date, err = service.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func addTaskHandler(srv *server.CustomServer) http.HandlerFunc {
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

		id, err := db.CreateTask(srv.Db, db.QueryInsertTask,
			task.Date,
			task.Title,
			task.Comment,
			task.Repeat)
		if err != nil {
			writeJson(w, err, http.StatusBadRequest, srv.Lg)
			return
		}

		srv.Lg.Printf("task %s succesfully created", task.Title)

		writeJson(w, miniID{ID: strconv.Itoa(int(id))}, http.StatusOK, srv.Lg)
	}
}
