package utils

import (
	"errors"
	"os"
	"time"

	"github.com/Bitummit/go_auth/internal/models"
	"github.com/Bitummit/go_auth/pkg/my_errors"
	"github.com/golang-jwt/jwt/v4"
)



type UserClaims struct {
	Id int64
	Username string
	ExpiresAt int64
}

func (u UserClaims) Valid() error {
	if u.ExpiresAt < time.Now().Unix() {
		return errors.New("token expired")
	}
	return nil
}


func NewToken(user models.User) (string, error) {
	duration, err := time.ParseDuration(os.Getenv("TOKEN_TTL"))
	if err != nil {
		return "", my_errors.ErrorTokenDuration
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		Id: user.Id,
		Username: user.Username,
		ExpiresAt: time.Now().Add(duration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", my_errors.ErrorSigningToken
	}

	return tokenString, err
}


func ParseToken(tokenString string) (models.User, error) {
	var userClaims UserClaims
	_, err := jwt.ParseWithClaims(tokenString, &userClaims, func(token *jwt.Token) (interface{}, error) {
    	return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return models.User{}, my_errors.ErrorInvalidToken
	}
	
	user := models.User{
		Id: userClaims.Id,
		Username: userClaims.Username,
	}

	return user, nil

}


// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MSwiVXNlcm5hbWUiOiJQZXJ0YXlhIiwiRXhwaXJlc0F0IjoxNzI4MzgzMzc5fQ.N1Vnsx0GHL6azLXeVfZwi3ik7W3dpBeTeKgv3nWlgkM