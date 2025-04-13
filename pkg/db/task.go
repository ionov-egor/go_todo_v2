package db

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func InsertTask(task *Task) (int, error) {
	res, err := GetDB().Exec("INSERT INTO scheduler  (date, title, comment, repeat) VALUES ( :date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastId), nil
}

func GetTaskById(id int) (task *Task, err error) {
	row := GetDB().QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", id))

	task = &Task{}

	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, err
	}

	return task, nil
}

func Tasks(limit int) ([]*Task, error) {
	var tasks = []*Task{}
	rows, err := GetDB().Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date DESC LIMIT :limit ", sql.Named("limit", limit))
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return tasks, err
		}

		tasks = append(tasks, &task)
	}

	err = rows.Err()
	if err != nil {
		return tasks, err
	}

	return tasks, nil
}

func TasksBySearch(limit int, search string) ([]*Task, error) {
	var tasks = []*Task{}

	rows, err := GetDB().Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE CONCAT('%', :search, '%')  OR comment LIKE CONCAT('%', :search, '%')  ORDER BY date DESC LIMIT :limit ",
		sql.Named("search", search),
		sql.Named("limit", limit))

	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return tasks, err
		}

		tasks = append(tasks, &task)
	}

	err = rows.Err()
	if err != nil {
		return tasks, err
	}

	return tasks, nil
}

func TasksByDate(limit int, date string) ([]*Task, error) {
	var tasks = []*Task{}

	rows, err := GetDB().Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date LIMIT :limit",
		sql.Named("date", date),
		sql.Named("limit", limit))

	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return tasks, err
		}

		tasks = append(tasks, &task)
	}

	err = rows.Err()
	if err != nil {
		return tasks, err
	}

	return tasks, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`
	res, err := GetDB().Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}

	return nil
}

func DeleteTask(id int) error {
	if _, err := GetDB().Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id)); err != nil {
		return err
	}

	return nil
}

func UpdateDate(next string, id int) error {
	query := `UPDATE scheduler SET date = :date WHERE id = :id`
	res, err := GetDB().Exec(query,
		sql.Named("date", next),
		sql.Named("id", id))
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}

	return nil

}
