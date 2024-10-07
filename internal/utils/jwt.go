package utils

import (
	"auth/internal/storage"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)



type UserClaims struct {
	Id int64
	Username string
	ExpiresAt int64
}

func (u UserClaims) Valid() error {
	if u.ExpiresAt < time.Now().Unix() {
		return errors.New("Token expired")
	}
	return nil
}


func NewToken(user storage.User) (string, error) {
	duration, err := time.ParseDuration(os.Getenv("TOKEN_TTL"))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		Id: user.Id,
		Username: user.Username,
		ExpiresAt: time.Now().Add(duration).Unix(),
	})
	
	// claims := token.Claims.(jwt.MapClaims)
	// claims["id"] = user.Id
	// claims["username"] = user.Username
	// claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", nil
	}

	return tokenString, err
}


func ParseToken(tokenString string) (storage.User, error) {
	var userClaims UserClaims
	_, err := jwt.ParseWithClaims(tokenString, &userClaims, func(token *jwt.Token) (interface{}, error) {
    	return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return storage.User{}, err
	}
	
	user := storage.User{
		Id: userClaims.Id,
		Username: userClaims.Username,
	}

	return user, nil

}


// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NSwiVXNlcm5hbWUiOiJQZXJ0YXlhIiwiRXhwaXJlc0F0IjoxNzI4MzQwNzM5fQ.z5SM4d_vAXcQ75kQSaC_DNYYDdsqL_fqRsUYq8WWQPM