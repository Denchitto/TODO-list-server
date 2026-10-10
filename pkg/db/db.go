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
	QueryById           = "SELECT * FROM scheduler WHERE id = ?;"
	QueryByDate         = "SELECT * FROM scheduler WHERE date = ? LIMIT ?;"
	QueryByTitleComment = "SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?;"
	QueryByAllResults   = "SELECT * FROM scheduler ORDER BY date LIMIT ?;"
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
		return nil, fmt.Errorf("error creating table: %w", err)
	}
	return db, nil
}
