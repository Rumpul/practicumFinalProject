package db

import (
	"errors"
	"time"

	"github.com/Yandex-Practicum/final-project/models"
	"github.com/Yandex-Practicum/final-project/utils"
	"github.com/jmoiron/sqlx"
)

func AddTask(db *sqlx.DB, task models.Task) (int64, error) {
	if task.Title == "" {
		return 0, errors.New("не указан заголовок задачи")
	}

	if task.Date == "" {
		task.Date = time.Now().Format(utils.TimeFormat)
	} else {
		_, err := time.Parse(utils.TimeFormat, task.Date)
		if err != nil {
			return 0, errors.New("дата представлена в неправильном формате")
		}
	}

	now := time.Now()
	if task.Date < now.Format(utils.TimeFormat) {
		if task.Repeat == "" {
			task.Date = now.Format(utils.TimeFormat)
		} else {
			nextDate, err := utils.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return 0, err
			}
			task.Date = nextDate
		}
	}

	result, err := db.Exec(
		`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func EditTask(db *sqlx.DB, task models.Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.Id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}

func GetTask(db *sqlx.DB, id string) (models.Task, error) {
	rows, err := db.NamedQuery(`SELECT id, date, title, comment, repeat 
	FROM scheduler WHERE id = :id`, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return models.Task{}, err
	}
	defer rows.Close()

	var res models.Task
	for rows.Next() {
		task := models.Task{}
		err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return res, err
		}
		res = task
	}
	err = rows.Err()
	if err != nil {
		return res, err
	}

	return res, nil
}

func GetTasks(db *sqlx.DB, limit int) ([]models.Task, error) {
	rows, err := db.NamedQuery(`SELECT id, date, title, comment, repeat 
	FROM scheduler ORDER BY date ASC LIMIT :limit;`, map[string]interface{}{
		"limit": limit,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.Task
	for rows.Next() {
		p := models.Task{}
		err := rows.Scan(&p.Id, &p.Date, &p.Title, &p.Comment, &p.Repeat)
		if err != nil {
			return nil, err
		}

		res = append(res, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func SearchTask(db *sqlx.DB, search string) ([]models.Task, error) {
	var rows *sqlx.Rows
	date, err := time.Parse("02.01.2006", search)
	if err != nil {
		rows, err = db.NamedQuery(`
		SELECT * FROM scheduler WHERE 
		id LIKE :input OR title LIKE :input
		OR comment LIKE :input ORDER BY date DESC;`,
			map[string]interface{}{
				"input": "%" + search + "%"})
	} else {
		rows, err = db.NamedQuery(`SELECT id, date, title, comment, repeat FROM scheduler 
		WHERE date = :searchDate ORDER BY date`,
			map[string]interface{}{
				"searchDate": date.Format(utils.TimeFormat),
			})
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.Task
	for rows.Next() {
		p := models.Task{}

		err := rows.Scan(&p.Id, &p.Date, &p.Title, &p.Comment, &p.Repeat)
		if err != nil {
			return nil, err
		}

		res = append(res, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
