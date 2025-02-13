package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	database "github.com/Yandex-Practicum/final-project/db"
	"github.com/Yandex-Practicum/final-project/models"
	"github.com/Yandex-Practicum/final-project/utils"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

const LimitTasks = 50

func HandleAddTask(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		var task models.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			log.Printf("ошибка десериализации JSON: %v", err)
			http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}

		id, err := database.AddTask(db, task)
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

func HandleEditTask(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			log.Printf("Ошибка десериализации JSON: %v", err)
			http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}

		if task.Id == "" {
			log.Printf("не указан идентификатор задачи")
			http.Error(w, `{"error": "Не указан идентификатор задачи"}`, http.StatusBadRequest)
			return
		}
		if task.Title == "" {
			log.Printf("не указан заголовок задачи")
			http.Error(w, `{"error": "Не указан заголовок задачи"}`, http.StatusBadRequest)
			return
		}
		if task.Date == "" {
			task.Date = time.Now().Format(utils.TimeFormat)
		} else {
			_, err := time.Parse(utils.TimeFormat, task.Date)
			if err != nil {
				log.Printf("дата представлена в неправильном формате")
				http.Error(w, `{"error": "Дата представлена в неправильном формате"}`, http.StatusBadRequest)
				return
			}
		}

		now := time.Now()
		if task.Date < now.Format(utils.TimeFormat) {
			if task.Repeat == "" {
				task.Date = now.Format(utils.TimeFormat)
			} else {
				nextDate, err := utils.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					log.Printf(`{"error": "` + err.Error() + `"}`)
					http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
					return
				}
				task.Date = nextDate
			}
		}
		err = database.EditTask(db, task)
		if err != nil {
			if err.Error() == "задача не найдена" {
				http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			} else {
				http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			}
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

func HandleGetTask(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		query := r.URL.Query()
		if !query.Has("id") {
			http.Error(w, `{"error": "Отсутствует индентификатор"}`, http.StatusInternalServerError)
			return
		}
		id := query.Get("id")
		var task models.Task
		task, err := database.GetTask(db, id)
		if err != nil {
			log.Printf("ошибка получения задач: %v", err)
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(task)
		if err != nil {
			log.Printf("не удалось закодировать ответ: %v", err)
		}
	}
}

func HandleGetTasks(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		search := r.URL.Query().Get("search")
		var tasks []models.Task
		var err error
		if len(search) > 0 {
			tasks, err = database.SearchTask(db, search)
		} else {
			tasks, err = database.GetTasks(db, LimitTasks)
		}
		if err != nil {
			log.Printf("ошибка получения задачь: %v", err)
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

func NextData(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Error string `json:"error,omitempty"`
		Date  string `json:"date,omitempty"`
	}{}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	query := r.URL.Query()
	params := []string{"now", "date", "repeat"}
	for _, param := range params {
		if !query.Has(param) {
			response.Error = fmt.Sprintf("пропущен обязательный параметр: %s", param)
			sendResponse(w, http.StatusBadRequest, response)
			return
		}
	}

	currDate, err := time.Parse(utils.TimeFormat, query.Get("now"))
	if err != nil {
		response.Error = fmt.Sprintf("неправильный формат даты: %v", err)
		sendResponse(w, http.StatusBadRequest, response)
		return
	}
	nextDate, err := utils.NextDate(
		currDate,
		query.Get("date"),
		query.Get("repeat"),
	)
	if err != nil {
		response.Error = err.Error()
		sendResponse(w, http.StatusBadRequest, response)
		return
	}
	response.Date = nextDate
	sendResponse(w, http.StatusOK, response)
}

func sendResponse(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("не удалось закодировать ответ: %v", err)
	}
}
