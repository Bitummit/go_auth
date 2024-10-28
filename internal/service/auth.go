package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Bitummit/go_auth/internal/models"
	"github.com/Bitummit/go_auth/internal/utils"
	"golang.org/x/crypto/bcrypt"
)


type( 
	AuthService struct {
		Storage UserStorage
	}

	UserStorage interface {
		CreateUser(ctx context.Context, user models.User) (int64, error)
		GetUser(ctx context.Context, username string) (*models.User, error)
	}
)

var ErrorWrongPassword = errors.New("wrong password")
var ErrorHashingPassword = errors.New("error while hashing password")


func New(storage UserStorage, log *slog.Logger) *AuthService {
	return &AuthService{
		Storage: storage,
	}
}

func (a *AuthService) CheckTokenUser(token string) error {
	user, err := utils.ParseToken(token)
	if err != nil {
		return fmt.Errorf("check user token: %w", err)
	}
	
	_, err = a.Storage.GetUser(context.Background(), user.Username)
	if err != nil {
		return fmt.Errorf("check user token: %w", err)
	}

	return nil
}


func (a *AuthService) LoginUser(username string, password string) (*string, error) {
	user, err := a.Storage.GetUser(context.Background(), username)
	if err != nil {
		return nil, fmt.Errorf("login user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password)); if err != nil {
		return nil, fmt.Errorf("login user: %w", err)
	}

	token, err := utils.NewToken(*user)
	if err != nil {
		return nil, fmt.Errorf("login user: %w", err)
	}

	return &token, nil
}


func (a *AuthService) RegisterUser(username string, password string) (*string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("generating password: %w", ErrorHashingPassword)
	}

	user := models.User{Username: username, Password: hashedPass}
	id, err := a.Storage.CreateUser(context.Background(), user)
	if err != nil {
		return nil, fmt.Errorf("registration user: %w", err)
	}

	user.Id = id
	token, err := utils.NewToken(user)
	if err != nil {
		return nil, fmt.Errorf("registration user: %w", err)
	}

	return &token, nil
}