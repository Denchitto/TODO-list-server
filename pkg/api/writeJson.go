package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeJson(w http.ResponseWriter, data any, status int, log *log.Logger) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var (
		marshalErr error
		msg        []byte
	)

	if err, ok := data.(error); ok {
		msg, marshalErr = json.Marshal(map[string]string{"error": err.Error()})
	} else {
		msg, marshalErr = json.Marshal(data)
	}

	if marshalErr != nil {
		log.Printf("marshal error: %v", marshalErr)
		http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(msg); err != nil {
		log.Printf("write error: %v", err)
		return
	}
}
