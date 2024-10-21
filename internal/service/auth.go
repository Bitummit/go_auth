package service

import (
	"context"
	"log/slog"

	"github.com/Bitummit/go_auth/internal/models"
	"github.com/Bitummit/go_auth/pkg/my_errors"
	"github.com/Bitummit/go_auth/internal/utils"
	"github.com/Bitummit/go_auth/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)


type UserStorage interface {
	CreateUser(context.Context, models.User) (int64, error)
	GetUser(context.Context, string) (*models.User, error)
}

type AuthService struct {
	storage UserStorage
	log *slog.Logger
}

func New(storage UserStorage, log *slog.Logger) *AuthService {
	return &AuthService{
		storage: storage,
		log: log,
	}
}

// TODO: переписать все ошибки на серверный слой
// Отсюда возвращать кастомную ошибку с %w чтобы сравнить ее в серверном слое
// Пробрасываем ошибку с самого низа


func (a *AuthService)CheckTokenUser(token string) error {
	user, err := utils.ParseToken(token)
	if err != nil {
		return err
	}
	
	_, err = a.storage.GetUser(context.Background(), user.Username)
	if err != nil {
		return err
	}

	return nil
}


func (a *AuthService)LoginUser(username string, password string) (*string, error) {
	user, err := a.storage.GetUser(context.Background(), username)
	if err != nil {
		// return nil, fmt.Errorf("error while fething data %v", err)
		a.log.Error("Getting user error: ", logger.Err(err))
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password)); if err != nil {
		// return nil, fmt.Errorf("wrong password %v", err)
		a.log.Error("Compare hash error: ", logger.Err(err))
		return nil, my_errors.ErrorWrongPassword
	}

	token, err := utils.NewToken(*user)
	if err != nil {
		// return nil, fmt.Errorf("error while generating token %v", err)
		a.log.Error("Generating token: ", logger.Err(err))
		return nil, err
	}
	return &token, nil
}


func (a *AuthService)RegisterUser(username string, password string) (*string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.log.Error("Hashing error: ", logger.Err(err))
		return nil, my_errors.ErrorHashingPassword
		// fmt.Errorf("error while hashing password %v", err)
	}

	user := models.User{Username: username, Password: hashedPass}
	id, err := a.storage.CreateUser(context.Background(), user)
	if err != nil {
		a.log.Error("User creation error: ", logger.Err(err))
		return nil, err
	}
	user.Id = id
	token, err := utils.NewToken(user)
	if err != nil {
		a.log.Error("Generating token error: ", logger.Err(err))
		return nil, err
	}

	return &token, nil
}