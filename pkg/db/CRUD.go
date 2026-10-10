package db

import (
	"database/sql"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func CreateTask(db *sql.DB, query string, args ...any) (int64, error) {

	res, err := db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func ReadTasks(db *sql.DB, query string, args ...any) ([]*Task, error) {

	rows, err := db.Query(query, args...)
	if err != nil {
		return []*Task{}, err
	}
	defer rows.Close()

	tasks := []*Task{}

	for rows.Next() {

		var newTask Task
		if err := rows.Scan(&newTask.ID,
			&newTask.Date,
			&newTask.Title,
			&newTask.Comment,
			&newTask.Repeat); err != nil {
			return []*Task{}, err
		}

		tasks = append(tasks, &newTask)
	}

	if err := rows.Err(); err != nil {
		return []*Task{}, err
	}
	return tasks, nil
}

func UpdateTask(db *sql.DB, query string, args ...any) (bool, error) {

	res, err := db.Exec(query, args...)
	if err != nil {
		return false, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return false, err
	} else if count == 0 {
		return true, nil
	}
	return false, nil
}

func DeleteTask(db *sql.DB, query string, args ...any) (bool, error) {

	res, err := db.Exec(query, args...)
	if err != nil {
		return false, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return false, err
	} else if count == 0 {
		return true, nil
	}
	return false, nil
}
