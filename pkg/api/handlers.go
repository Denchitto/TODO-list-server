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

func nextDayHandler(w http.ResponseWriter, r *http.Request) {

	keys := []string{"now", "date", "repeat"}
	for i, v := range keys {
		keys[i] = r.FormValue(v)
	}
	if keys[0] == "" {
		keys[0] = time.Now().Format(service.DateFormat)
	}

	now, err := time.Parse(service.DateFormat, keys[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data, err := service.NextDate(now, keys[1], keys[2])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}

func writeJson(w http.ResponseWriter, msg string, status int) error {
	w.Header().Set("Content-Type", "application/json; charser=UTF-8")
	w.WriteHeader(status)
	jsonForm := make(map[string]string)
	switch {
	case status == http.StatusOK:
		jsonForm["id"] = msg
	case status != http.StatusOK:
		jsonForm["error"] = msg
	}
	data, err := json.Marshal(jsonForm)
	if err != nil {
		return err
	}
	w.Write(data)
	return nil
}

func addTaskHandler(srv *server.CustomServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			task       db.Task
			buf        bytes.Buffer
			err        error
			parsedDate time.Time
		)

		_, err = buf.ReadFrom(r.Body)
		if err != nil {
			if err = writeJson(w, err.Error(), http.StatusBadRequest); err != nil {
				srv.Lg.Fatal("marshal error: %w", err)
			}
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			if err = writeJson(w, err.Error(), http.StatusBadRequest); err != nil {
				srv.Lg.Fatal("marshal error: %w", err)
			}
			return
		}

		if task.Title == "" {
			if err = writeJson(w, fmt.Errorf("empty title").Error(), http.StatusBadRequest); err != nil {
				srv.Lg.Fatal("marshal error: %w", err)
			}
			return
		}

		now := time.Now()
		switch {
		case strings.TrimSpace(task.Date) == "":
			task.Date = now.Format(service.DateFormat)
		default:
			parsedDate, err = time.Parse(service.DateFormat, task.Date)
			if err != nil {
				if err = writeJson(w, err.Error(), http.StatusBadRequest); err != nil {
					srv.Lg.Fatal("marshal error: %w", err)
				}
				return
			}
			if service.AfterNow(now, parsedDate) {
				switch {
				case strings.TrimSpace(task.Repeat) == "":
					task.Date = now.Format(service.DateFormat)
				default:
					task.Date, err = service.NextDate(now, task.Date, task.Repeat)
					if err != nil {
						if err = writeJson(w, err.Error(), http.StatusBadRequest); err != nil {
							srv.Lg.Fatal("marshal error: %w", err)
						}
						return
					}
				}
			}
		}

		id, err := db.AddTask(srv.Db, &task)
		if err != nil {
			if err = writeJson(w, err.Error(), http.StatusBadRequest); err != nil {
				srv.Lg.Fatal("marshal error: %w", err)
			}
			return
		}

		if err = writeJson(w, strconv.Itoa(int(id)), http.StatusOK); err != nil {
			srv.Lg.Fatal("marshal error: %w", err)
		}
	}
}
