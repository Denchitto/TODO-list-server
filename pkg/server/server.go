package server

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"
)

type CustomServer struct {
	Lg  *log.Logger
	Srv *http.Server
	Db  *sql.DB
}

func NewServer(port string) (*CustomServer, error) {

	customSrv := &CustomServer{}

	customSrv.Lg = log.New(os.Stdout, "INFO: ", log.LstdFlags)

	customSrv.Srv = &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	return customSrv, nil
}
