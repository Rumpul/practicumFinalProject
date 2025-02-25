package service

import (
	"errors"
	"time"

	"github.com/Yandex-Practicum/final-project/dates"
	"github.com/Yandex-Practicum/final-project/models"
	"github.com/Yandex-Practicum/final-project/storage"
)

type TaskService struct {
	storage *storage.TaskStorage
}

func NewTaskService(storage *storage.TaskStorage) *TaskService {
	return &TaskService{storage: storage}
}

func (s *TaskService) AddTask(task models.Task) (int64, error) {
	if task.Title == "" {
		return 0, errors.New("не указан заголовок задачи")
	}

	if task.Date == "" {
		task.Date = time.Now().Format(dates.TimeFormat)
	} else {
		if _, err := time.Parse(dates.TimeFormat, task.Date); err != nil {
			return 0, errors.New("неправильный формат даты")
		}
	}

	now := time.Now()
	currentDate := now.Format(dates.TimeFormat)

	if task.Date < currentDate {
		if task.Repeat == "" {
			task.Date = currentDate
		} else {
			nextDate, err := dates.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return 0, err
			}
			task.Date = nextDate
		}
	}

	return s.storage.AddTask(task)
}

func (s *TaskService) EditTask(task models.Task) error {
	if task.Title == "" {
		return errors.New("заголовок не может быть пустым")
	}
	return s.storage.EditTask(task)
}

func (s *TaskService) GetTask(id string) (models.Task, error) {
	return s.storage.GetTask(id)
}

func (s *TaskService) GetTasks(limit int) ([]models.Task, error) {
	return s.storage.GetTasks(limit)
}

func (s *TaskService) SearchTasks(search string) ([]models.Task, error) {
	date, err := time.Parse("02.01.2006", search)
	if err == nil {
		return s.storage.SearchByDate(date.Format(dates.TimeFormat))
	}
	return s.storage.SearchByText(search)
}

func (s *TaskService) DeleteTask(id string) error {
	return s.storage.DeleteTask(id)
}
