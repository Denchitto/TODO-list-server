package db

import (
	"database/sql"
	"strings"
	"time"

	"github.com/Denchitto/TODO-list-server/pkg/service"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func CreateTask(db *sql.DB, args ...any) (int64, error) {

	res, err := db.Exec(queryInsertTask, args...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func ReadTask(db *sql.DB, args ...any) (*Task, error) {

	var newTask Task
	if err := db.QueryRow(queryById, args...).Scan(
		&newTask.ID,
		&newTask.Date,
		&newTask.Title,
		&newTask.Comment,
		&newTask.Repeat); err != nil {
		return &Task{}, err
	}
	return &newTask, nil
}

func ReadTasks(db *sql.DB, value string, limit int) ([]*Task, error) {

	query, args := readQueryArgs(value, limit)

	rows, err := db.Query(query, args...)
	if err != nil {
		return []*Task{}, err
	}
	defer rows.Close()

	tasks := []*Task{}

	for rows.Next() {

		var newTask Task
		if err := rows.Scan(
			&newTask.ID,
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

func readQueryArgs(value string, limit int) (string, []any) {

	var (
		query string
		args  []any
	)

	if strings.TrimSpace(value) == "" {
		query = queryByAllResults
		args = []any{limit}
	} else if t, err := time.Parse(service.DateFormatSearch, value); err == nil {
		value = t.Format(service.DateFormat)
		query = queryByDate
		args = []any{value, limit}
	} else {
		value = "%" + value + "%"
		query = queryByTitleComment
		args = []any{value, value, limit}
	}
	return query, args
}

func UpdateTask(db *sql.DB, args ...any) (int64, error) {

	res, err := db.Exec(queryUpdateTask, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func UpdateDate(db *sql.DB, args ...any) (int64, error) {

	res, err := db.Exec(queryUpdateTaskByDate, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func DeleteTask(db *sql.DB, args ...any) (int64, error) {

	res, err := db.Exec(queryDeleteTask, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
