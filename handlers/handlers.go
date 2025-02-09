package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Yandex-Practicum/final-project/utils"
)

func NextData(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	now := r.URL.Query().Get("now")
	currDate, err := time.Parse(utils.TimeFormat, now)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")
	nextdate, err := utils.NextDate(currDate, date, repeat)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	answer, err := strconv.Atoi(nextdate)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	err = json.NewEncoder(w).Encode(answer)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
}
