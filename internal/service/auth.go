package service

import (
	"context"
	"fmt"

	"github.com/Bitummit/go_auth/internal/models"
	"github.com/Bitummit/go_auth/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/internal/status"
)


type UserStorage interface {
	CreateUser(context.Context, models.User) (int64, error)
	GetUser(context.Context, string) (*models.User, error)
}

type AuthService struct {
	storage UserStorage
}

func New(storage UserStorage) *AuthService {
	return &AuthService{
		storage: storage,
	}
}

// TODO: переписать все ошибки на серверный слой
// Отсюда возвращать кастомную ошибку с %w чтобы сравнить ее в серверном слое



func (a *AuthService)CheckTokenUser(token string) error {
	user, err := utils.ParseToken(token)
	if err != nil {
		return status.Err(codes.InvalidArgument, err)
	}
	
	_, err = a.storage.GetUser(context.Background(), user.Username)
	if err != nil {
		return status.Err(codes.NotFound, err)
	}

	return nil
}


func (a *AuthService)LoginUser(username string, password string) (*string, error) {
	user, err := a.storage.GetUser(context.Background(), username)
	if err != nil {
		// return nil, fmt.Errorf("error while fething data %v", err)
		return nil, status.Err(codes.NotFound, err)
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password)); if err != nil {
		// return nil, fmt.Errorf("wrong password %v", err)
		return nil, status.Err(codes.InvalidArgument, "wrong password")
	}

	token, err := utils.NewToken(*user)
	if err != nil {
		// return nil, fmt.Errorf("error while generating token %v", err)
		return nil, status.Err(codes.InvalidArgument, err)
	}
	return &token, nil
}


func (a *AuthService)RegisterUser(username string, password string) (*string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Err(codes.InvalidArgument, "err")
		fmt.Errorf("error while hashing password %v", err)
	}

	user := models.User{Username: username, Password: hashedPass}
	id, err := a.storage.CreateUser(context.Background(), user)
	if err != nil {
		return nil, fmt.Errorf("error while inserting user %v", err)
	}
	user.Id = id
	token, err := utils.NewToken(user)
	if err != nil {
		return nil, fmt.Errorf("error while generating token %v", err)
	}

	return &token, nil
}