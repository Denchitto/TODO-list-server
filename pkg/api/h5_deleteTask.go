package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Denchitto/TODO-list-server/pkg/db"
)

func deleteTaskHandler(database *sql.DB, lg *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Берем id и проверяем его
		id := r.URL.Query().Get("id")

		if strings.TrimSpace(id) == "" {
			writeJson(w, fmt.Errorf("id not specified"), http.StatusBadRequest, lg)
			return
		}
		if _, err := strconv.Atoi(id); err != nil {
			writeJson(w, fmt.Errorf("incorrect id format: %s; %w", id, err), http.StatusBadRequest, lg)
			return
		}
		// 2. Удаляем запись, проверяем что запись удалилась
		count, err := db.DeleteTask(database, id)

		if err != nil {
			lg.Println(err)
			writeJson(w, fmt.Errorf("internal error"), http.StatusInternalServerError, lg)
			return
		} else if count == 0 {
			writeJson(w, fmt.Errorf("incorrect id for deleting task"), http.StatusNotFound, lg)
			return
		}

		lg.Println("task succesfully deleted")

		writeJson(w, struct{}{}, http.StatusOK, lg)
	}
}
