package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/final-project/dates"
	"github.com/Yandex-Practicum/final-project/models"
	"github.com/Yandex-Practicum/final-project/service"
)

const LimitTasks = 50

func HandleAddTask(service *service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		var task models.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			log.Printf("ошибка десериализации JSON: %v", err)
			http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}

		id, err := service.AddTask(task)
		if err != nil {
			log.Printf("ошибка добавления задачи: %v", err)
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
		if err != nil {
			log.Printf("не удалось закодировать ответ: %v", err)
		}
	}
}

func HandleEditTask(service *service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		var task models.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			log.Printf("ошибка десериализации JSON: %v", err)
			http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}

		if task.Id == "" {
			log.Println("не указан идентификатор задачи")
			http.Error(w, `{"error": "Не указан идентификатор задачи"}`, http.StatusBadRequest)
			return
		}
		if task.Title == "" {
			log.Println("не указан заголовок задачи")
			http.Error(w, `{"error": "Не указан заголовок задачи"}`, http.StatusBadRequest)
			return
		}

		if task.Date == "" {
			task.Date = time.Now().Format(dates.TimeFormat)
		} else {
			_, err := time.Parse(dates.TimeFormat, task.Date)
			if err != nil {
				log.Printf("дата представлена в неправильном формате")
				http.Error(w, `{"error": "Дата представлена в неправильном формате"}`, http.StatusBadRequest)
				return
			}
		}

		now := time.Now()
		if task.Date < now.Format(dates.TimeFormat) {
			if task.Repeat == "" {
				task.Date = now.Format(dates.TimeFormat)
			} else {
				nextDate, err := dates.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					log.Println(err.Error())
					http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
					return
				}
				task.Date = nextDate
			}
		}

		err = service.EditTask(task)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(map[string]interface{}{})
		if err != nil {
			log.Printf("не удалось закодировать ответ: %v", err)
		}
	}
}

func HandleGetTask(service *service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		query := r.URL.Query()

		if !query.Has("id") {
			log.Println("отсутствует идентификатор")
			http.Error(w, `{"error": "Отсутствует идентификатор"}`, http.StatusBadRequest)
			return
		}

		id := query.Get("id")
		task, err := service.GetTask(id)
		if err != nil {
			log.Printf("ошибка получения задачи: %v", err)
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(task)
		if err != nil {
			log.Printf("не удалось закодировать ответ: %v", err)
		}
	}
}

func HandleGetTasks(service *service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		search := r.URL.Query().Get("search")
		var tasks []models.Task
		var err error

		if search != "" {
			tasks, err = service.SearchTasks(search)
		} else {
			tasks, err = service.GetTasks(LimitTasks)
		}

		if err != nil {
			log.Printf("ошибка получения задач: %v", err)
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		if tasks == nil {
			tasks = []models.Task{}
		}

		err = json.NewEncoder(w).Encode(map[string]interface{}{"tasks": tasks})
		if err != nil {
			log.Printf("не удалось закодировать ответ: %v", err)
		}
	}
}

func HandleDeleteTask(service *service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		query := r.URL.Query()

		if !query.Has("id") {
			log.Println("отсутствует идентификатор")
			http.Error(w, `{"error": "Отсутствует идентификатор"}`, http.StatusBadRequest)
			return
		}

		id := query.Get("id")
		err := service.DeleteTask(id)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(map[string]interface{}{})
		if err != nil {
			log.Printf("не удалось закодировать ответ: %v", err)
		}
	}
}

func HandleTaskDone(service *service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		query := r.URL.Query()

		if !query.Has("id") {
			log.Println("отсутствует идентификатор")
			http.Error(w, `{"error": "Отсутствует идентификатор"}`, http.StatusBadRequest)
			return
		}

		id := query.Get("id")
		task, err := service.GetTask(id)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		if task.Repeat == "" {
			err := service.DeleteTask(id)
			if err != nil {
				http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
				return
			}
		} else {
			now := time.Now()
			nextDate, err := dates.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, `{"error": "Ошибка вычисления следующей даты"}`, http.StatusInternalServerError)
				return
			}
			task.Date = nextDate
			err = service.EditTask(task)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
				return
			}
		}

		err = json.NewEncoder(w).Encode(map[string]interface{}{})
		if err != nil {
			log.Printf("не удалось закодировать ответ: %v", err)
		}
	}
}

func NextData(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	params := []string{"now", "date", "repeat"}
	for _, param := range params {
		if !query.Has(param) {
			log.Printf("пропущен обязательный параметр: %s", param)
			http.Error(w, `{"error": "Пропущен обязательный параметр"}`, http.StatusBadRequest)
			return
		}
	}

	currDate, err := time.Parse(dates.TimeFormat, query.Get("now"))
	if err != nil {
		log.Printf("неправильный формат даты: %v", err)
		http.Error(w, `{"error": "Неправильный формат даты"}`, http.StatusBadRequest)
		return
	}
	nextDate, err := dates.NextDate(
		currDate,
		query.Get("date"),
		query.Get("repeat"),
	)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	w.Write([]byte(nextDate))
}
