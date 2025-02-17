package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/Yandex-Practicum/final-project/jwt"
	"github.com/Yandex-Practicum/final-project/models"
)

var password = os.Getenv("TODO_PASSWORD")

func HangdleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if len(password) > 0 {
		var login models.Login
		err := json.NewDecoder(r.Body).Decode(&login)
		if err != nil {
			log.Printf("ошибка десериализации JSON: %v", err)
			http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}
		if login.Password != password {
			log.Println("некорректные данные")
			http.Error(w, `{"error": "Некорректные данные"}`, http.StatusForbidden)
			return
		}
		newToken, err := jwt.JWTCreate()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(newToken)
		if err != nil {
			log.Printf("не удалось закодировать ответ: %v", err)
		}
	}
}
