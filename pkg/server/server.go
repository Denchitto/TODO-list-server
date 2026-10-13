package server

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Denchitto/TODO-list-server/pkg/api"
	"github.com/Denchitto/TODO-list-server/pkg/db"
)

type CustomServer struct {
	Lg  *log.Logger
	Db  *sql.DB
	Srv *http.Server
}

func NewServer(port, path string) (*CustomServer, error) {

	lg := log.New(os.Stdout, "INFO: ", log.LstdFlags)

	dbConn, err := db.DbConnect(path)
	if err != nil {
		return nil, err
	}
	lg.Println("Init DB")

	mux, err := api.Init(dbConn, lg)
	if err != nil {
		dbConn.Close()
		return nil, err
	}

	lg.Printf("Init server on port %s", port)

	return &CustomServer{
		Lg: lg,
		Db: dbConn,
		Srv: &http.Server{
			Addr:         ":" + port,
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  10 * time.Second,
		},
	}, nil
}
