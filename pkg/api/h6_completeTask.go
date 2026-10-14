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

func completeTaskHandler(database *sql.DB, lg *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Извлекаем id и проверяем его
		id := r.URL.Query().Get("id")
		if strings.TrimSpace(id) == "" {
			writeJson(w, fmt.Errorf("id not specified"), http.StatusBadRequest, lg)
			return
		}
		if _, err := strconv.Atoi(id); err != nil {
			writeJson(w, fmt.Errorf("incorrect id format: %s; %w", id, err), http.StatusBadRequest, lg)
			return
		}
		// 2. Ищем задачу по id, проверяем, что id найден, других ошибок нет
		task, err := db.ReadTask(database, db.QueryById, id)
		if errors.Is(sql.ErrNoRows, err) {
			writeJson(w, fmt.Errorf("unable to find task with id: %s", id), http.StatusNotFound, lg)
			return
		}
		if err != nil {
			writeJson(w, err, http.StatusInternalServerError, lg)
			return
		}
		// 3. Удаляем задачу если повтора нет
		if strings.TrimSpace(task.Repeat) == "" {
			if _, err := db.DeleteTask(database, db.QueryDeleteTask, task.ID); err != nil {
				writeJson(w, err, http.StatusInternalServerError, lg)
				return
			}

			lg.Printf("task %s succesfully deleted", task.Title)

			writeJson(w, struct{}{}, http.StatusOK, lg)
			return
		} else { // Иначе вычисляем следующую дату и в базу ее
			task.Date, err = service.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				writeJson(w, err, http.StatusInternalServerError, lg)
				return
			}
			if _, err = db.UpdateTask(database, db.QueryUpdateTaskById, task.Date, task.ID); err != nil {
				writeJson(w, err, http.StatusInternalServerError, lg)
				return
			}

			lg.Printf("task %s succesfully updated", task.Title)

			writeJson(w, struct{}{}, http.StatusOK, lg)
		}
	}
}
