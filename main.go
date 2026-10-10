package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Denchitto/TODO-list-server/pkg/api"
	"github.com/Denchitto/TODO-list-server/pkg/db"
	"github.com/Denchitto/TODO-list-server/pkg/server"
)

func getPort() (string, error) {

	port := os.Getenv("TODO_PORT")
	if strings.TrimSpace(port) == "" {
		log.Println("Setting default port: 7540")
		return "7540", nil
	}

	num, err := strconv.Atoi(port)
	if err != nil {
		return "", fmt.Errorf("incorrect port format: %w", err)
	}

	if num < 0 || num > 65535 {
		return "", fmt.Errorf("port value %d is out of range", num)
	}
	return strconv.Itoa(num), nil
}

func getPath(dbName string) string {

	path := os.Getenv("TODO_DBFILE")
	if strings.TrimSpace(path) == "" {
		log.Println("Setting current directory for DB as default")
		return dbName
	}

	return filepath.Join(filepath.Clean(path), dbName)
}

func main() {

	port, err := getPort()
	if err != nil {
		log.Fatal(err)
	}

	path := getPath("scheduler.db")

	srv, err := server.NewServer(port)
	if err != nil {
		log.Fatal(err)
	}
	srv.Lg.Printf("Init server on port %s", srv.Srv.Addr)

	dbConn, err := db.DbConnect(path)
	if err != nil {
		srv.Lg.Fatal(err)
	}
	srv.Db = dbConn
	defer srv.Db.Close()
	srv.Lg.Println("Init DB")

	mux, err := api.NewRouter()
	if err != nil {
		srv.Lg.Fatal(err)
	}

	api.Init(mux, srv)
	srv.Srv.Handler = mux

	if err = srv.Srv.ListenAndServe(); err != nil {
		srv.Lg.Println(err)
	}
	srv.Lg.Println("Have a nice day!")
}
