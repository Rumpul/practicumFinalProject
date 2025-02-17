package storage

import (
	"errors"
	"time"

	"github.com/Yandex-Practicum/final-project/dates"
	"github.com/Yandex-Practicum/final-project/models"
	"github.com/jmoiron/sqlx"
)

type TaskStorage struct {
	db *sqlx.DB
}

func NewTaskStorage(db *sqlx.DB) *TaskStorage {
	return &TaskStorage{db: db}
}

func (s *TaskStorage) AddTask(task models.Task) (int64, error) {
	if task.Title == "" {
		return 0, errors.New("не указан заголовок задачи")
	}

	if task.Date == "" {
		task.Date = time.Now().Format(dates.TimeFormat)
	} else {
		_, err := time.Parse(dates.TimeFormat, task.Date)
		if err != nil {
			return 0, errors.New("дата представлена в неправильном формате")
		}
	}

	now := time.Now()
	if task.Date < now.Format(dates.TimeFormat) {
		if task.Repeat == "" {
			task.Date = now.Format(dates.TimeFormat)
		} else {
			nextDate, err := dates.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return 0, err
			}
			task.Date = nextDate
		}
	}

	result, err := s.db.Exec(
		`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (s *TaskStorage) EditTask(task models.Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := s.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.Id)
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

func (s *TaskStorage) GetTask(id string) (models.Task, error) {
	var task models.Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := s.db.QueryRow(query, id).Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, err
	}
	return task, nil
}

func (s *TaskStorage) GetTasks(limit int) ([]models.Task, error) {
	rows, err := s.db.NamedQuery(`SELECT id, date, title, comment, repeat 
	FROM scheduler ORDER BY date DESC LIMIT :limit;`, map[string]interface{}{
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

	return res, rows.Err()
}

func (s *TaskStorage) SearchTask(search string) ([]models.Task, error) {
	var rows *sqlx.Rows
	date, err := time.Parse("02.01.2006", search)
	if err != nil {
		rows, err = s.db.NamedQuery(`
		SELECT * FROM scheduler WHERE 
		id LIKE :input OR title LIKE :input
		OR comment LIKE :input ORDER BY date DESC LIMIT :limit;`,
			map[string]interface{}{
				"input": "%" + search + "%",
				"limit": 100, // Добавляем лимит по умолчанию
			})
	} else {
		rows, err = s.db.NamedQuery(`SELECT id, date, title, comment, repeat FROM scheduler 
		WHERE date = :searchDate ORDER BY date DESC LIMIT :limit`,
			map[string]interface{}{
				"searchDate": date.Format(dates.TimeFormat),
				"limit":      100,
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

	return res, rows.Err()
}

func (s *TaskStorage) DeleteTask(id string) error {
	res, err := s.db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
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

func (s *TaskStorage) Close() error {
	return s.db.Close()
}
