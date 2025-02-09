package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/final-project/utils"
)

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
