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

	QueryInsertTask     = "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?);"
	QueryById           = "SELECT s.id, s.date, s.title, s.comment, s.repeat FROM scheduler s WHERE s.id = ?;"
	QueryByDate         = "SELECT s.id, s.date, s.title, s.comment, s.repeat FROM scheduler s WHERE s.date = ? LIMIT ?;"
	QueryByTitleComment = "SELECT s.id, s.date, s.title, s.comment, s.repeat FROM scheduler s WHERE s.title LIKE ? OR s.comment LIKE ? ORDER BY s.date LIMIT ?;"
	QueryByAllResults   = "SELECT s.id, s.date, s.title, s.comment, s.repeat FROM scheduler s ORDER BY s.date LIMIT ?;"
	QueryUpdateTask     = "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?;"
	QueryUpdateTaskById = "UPDATE scheduler SET date = ? WHERE id = ?;"
	QueryDeleteTask     = "DELETE FROM scheduler WHERE id = ?;"
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
