package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const (
	schema = `CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title TEXT NOT NULL DEFAULT "",
		comment TEXT NOT NULL DEFAULT "",
		repeat VARCHAR(128) NOT NULL DEFAULT "" 
	);
	CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);`

	queryInsertTask       = "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?);"
	queryById             = "SELECT s.id, s.date, s.title, s.comment, s.repeat FROM scheduler s WHERE s.id = ?;"
	queryByDate           = "SELECT s.id, s.date, s.title, s.comment, s.repeat FROM scheduler s WHERE s.date = ? LIMIT ?;"
	queryByTitleComment   = "SELECT s.id, s.date, s.title, s.comment, s.repeat FROM scheduler s WHERE s.title LIKE ? OR s.comment LIKE ? ORDER BY s.date LIMIT ?;"
	queryByAllResults     = "SELECT s.id, s.date, s.title, s.comment, s.repeat FROM scheduler s ORDER BY s.date LIMIT ?;"
	queryUpdateTask       = "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?;"
	queryUpdateTaskByDate = "UPDATE scheduler SET date = ? WHERE id = ?;"
	queryDeleteTask       = "DELETE FROM scheduler WHERE id = ?;"
)

func DbConnect(dbFile string) (*sql.DB, error) {

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("unable to establish the connection with database: %w", err)
	}

	if _, err = db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("error creating table: %w", err)
	}
	return db, nil
}
