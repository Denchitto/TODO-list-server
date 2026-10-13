package api

import (
	"log"
	"net/http"
	"time"

	"github.com/Denchitto/TODO-list-server/pkg/service"
)

func nextDateHandler(lg *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		keys := []string{"now", "date", "repeat"}
		for i, v := range keys {
			keys[i] = r.FormValue(v)
		}
		if keys[0] == "" {
			keys[0] = time.Now().Format(service.DateFormat)
			lg.Println("'now' not provided, setting current time")
		}

		now, err := time.Parse(service.DateFormat, keys[0])
		if err != nil {
			lg.Printf("error parsing 'now': %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		data, err := service.NextDate(now, keys[1], keys[2])
		if err != nil {
			lg.Printf("error calculating next date: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))
	}
}
