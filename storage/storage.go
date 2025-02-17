package storage

import (
	"database/sql"
	"errors"

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
	err := s.db.QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`,
		id,
	).Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err == sql.ErrNoRows {
		return task, errors.New("задача не найдена")
	}
	return task, err
}

func (s *TaskStorage) GetTasks(limit int) ([]models.Task, error) {
	var tasks []models.Task
	err := s.db.Select(&tasks,
		`SELECT id, date, title, comment, repeat 
		FROM scheduler ORDER BY date DESC LIMIT ?`,
		limit,
	)
	return tasks, err
}

func (s *TaskStorage) SearchByDate(date string) ([]models.Task, error) {
	var tasks []models.Task
	err := s.db.Select(&tasks,
		`SELECT id, date, title, comment, repeat 
		FROM scheduler WHERE date = ? ORDER BY date DESC`,
		date,
	)
	return tasks, err
}

func (s *TaskStorage) SearchByText(text string) ([]models.Task, error) {
	var tasks []models.Task
	err := s.db.Select(&tasks,
		`SELECT id, date, title, comment, repeat 
		FROM scheduler 
		WHERE title LIKE ? OR comment LIKE ? 
		ORDER BY date DESC LIMIT 100`,
		"%"+text+"%", "%"+text+"%",
	)
	return tasks, err
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
