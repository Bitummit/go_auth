package utils

import (
	"auth/internal/storage"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)


func NewToken(user storage.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	duration, err := time.ParseDuration(os.Getenv("TOKEN_TTL"))
	if err != nil {
		return "", err
	}
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.Id
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", nil
	}

	return tokenString, err
}