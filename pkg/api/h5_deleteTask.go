package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Denchitto/TODO-list-server/pkg/db"
	"github.com/Denchitto/TODO-list-server/pkg/server"
)

func deleteTaskHandler(srv *server.CustomServer) http.HandlerFunc {
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

		count, err := db.DeleteTask(srv.Db, db.QueryDeleteTask, id)

		if err != nil {
			writeJson(w, err, http.StatusInternalServerError, srv.Lg)
			return
		} else if count == 0 {
			writeJson(w, fmt.Errorf("incorrect id for deleting task"), http.StatusNotFound, srv.Lg)
			return
		}

		srv.Lg.Println("task succesfully deleted")

		writeJson(w, struct{}{}, http.StatusOK, srv.Lg)
	}
}
