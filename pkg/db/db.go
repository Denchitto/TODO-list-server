package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const (
	schema = `CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8),
		title TEXT NOT NULL DEFAULT "",
		comment TEXT NOT NULL DEFAULT "",
		repeat VARCHAR(128) 
	);
	CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);`
)

func DbConnect(dbFile string) (*sql.DB, error) {
	var err error

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("unable to establish the connection with database: %w", err)
	}

	if _, err = db.Exec(schema); err != nil {
		return nil, fmt.Errorf("wew %w", err)
	}
	return db, nil
}

func AddTask(db *sql.DB, task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES 
				(:date, :title, :comment, :repeat);`
	res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}
	id, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
