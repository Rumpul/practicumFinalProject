package middleware

import (
	"net/http"
	"os"

	"github.com/Yandex-Practicum/final-project/jwt"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var cookieToken string
			cookie, err := r.Cookie("token")
			if err == nil {
				cookieToken = cookie.Value
			}
			valid := jwt.JWTValidate(cookieToken)

			if valid != nil {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
