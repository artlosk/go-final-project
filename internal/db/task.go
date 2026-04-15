package db

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	result, err := DB.Exec(
		"INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err == nil {
		id, err = result.LastInsertId()
	}
	return id, err
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`неверный id`)
	}
	return nil
}

func UpdateDate(next string, id string) error {
	res, err := DB.Exec(`UPDATE scheduler SET date=? WHERE id=?`, next, id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func DeleteTask(id string) error {
	res, err := DB.Exec(`DELETE FROM scheduler WHERE id=?`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func GetTask(id string) (*Task, error) {
	task := &Task{}
	var parsedID int64

	err := DB.QueryRow(
		"SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?",
		id,
	).Scan(&parsedID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}

	task.ID = strconv.FormatInt(parsedID, 10)
	return task, nil
}

func Tasks(limit int, search string) ([]*Task, error) {
	tasks := make([]*Task, 0)

	query := `
SELECT id, date, title, comment, repeat
FROM scheduler
ORDER BY date
LIMIT ?
`
	args := []any{limit}
	if search != "" {
		query = `
SELECT id, date, title, comment, repeat
FROM scheduler
WHERE title LIKE ? COLLATE NOCASE
   OR comment LIKE ? COLLATE NOCASE
   OR date = ?
ORDER BY date
LIMIT ?
`
		like := "%" + search + "%"
		dateQuery := ""
		if parsedDate, err := time.Parse("02.01.2006", search); err == nil {
			dateQuery = parsedDate.Format("20060102")
		}
		args = []any{like, like, dateQuery, limit}
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		return tasks, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("ошибка: %v", err)
		}
	}(rows)

	for rows.Next() {
		task := Task{}
		var id int64

		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return tasks, err
		}
		task.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return tasks, err
	}

	return tasks, nil
}
