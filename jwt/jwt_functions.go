package jwt

import (
	"errors"
	"os"

	"github.com/Yandex-Practicum/final-project/models"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = os.Getenv("SECRET_KEY")

func JWTCreate() (models.LoginResponse, error) {
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	newToken, err := jwtToken.SignedString([]byte(secretKey))
	if err != nil {
		return models.LoginResponse{}, err
	}
	return models.LoginResponse{Token: newToken}, nil
}

func JWTValidate(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный метод подписи")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("токен недействителен")
	}
	return nil
}
