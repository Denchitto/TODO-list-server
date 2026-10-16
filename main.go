package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Denchitto/TODO-list-server/pkg/server"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	port, err := getPort()
	if err != nil {
		log.Fatal(err)
	}

	path := getPath()

	srv, err := server.NewServer(port, path)
	if err != nil {
		log.Fatal(err)
	}
	defer srv.Db.Close()

	if err = srv.Srv.ListenAndServe(); err != nil {
		srv.Lg.Println(err)
	}
}

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
	// Проверяем что порт в пределах допустимого диапазона портов
	if num < 0 || num > 65535 {
		return "", fmt.Errorf("port value %d is out of range", num)
	}
	return port, nil
}

func getPath() string {

	path := os.Getenv("TODO_DBFILE")
	if strings.TrimSpace(path) == "" {
		log.Println("Setting current directory for DB as default")
		return "scheduler.db"
	}
	return filepath.Clean(path)
}
