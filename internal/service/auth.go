package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	my_kafka "github.com/Bitummit/go_auth/internal/api/kafka"
	my_jwt "github.com/Bitummit/go_auth/internal/jwt"
	"github.com/Bitummit/go_auth/internal/models"
	"golang.org/x/crypto/bcrypt"
)


type( 
	AuthService struct {
		Storage UserStorage
		Kafka *my_kafka.Kafka
	}

	UserStorage interface {
		CreateUser(ctx context.Context, user models.User) (int64, error)
		GetUser(ctx context.Context, username string) (*models.User, error)
	}
)

var ErrorWrongPassword = errors.New("wrong password")
var ErrorHashingPassword = errors.New("error while hashing password")


func New(storage UserStorage, log *slog.Logger, kafka *my_kafka.Kafka) *AuthService {
	return &AuthService{
		Storage: storage,
		Kafka: kafka,
	}
}

func (a *AuthService) CheckTokenUser(_ context.Context, token string) error {
	user, err := my_jwt.ParseToken(token)
	if err != nil {
		return fmt.Errorf("check user token: %w", err)
	}
	
	_, err = a.Storage.GetUser(context.Background(), user.Username)
	if err != nil {
		return fmt.Errorf("check user token: %w", err)
	}

	return nil
}


func (a *AuthService) LoginUser(_ context.Context, username string, password string) (*string, error) {
	user, err := a.Storage.GetUser(context.Background(), username)
	if err != nil {
		return nil, fmt.Errorf("login user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password)); if err != nil {
		return nil, fmt.Errorf("login user: %w", err)
	}

	token, err := my_jwt.NewToken(*user)
	if err != nil {
		return nil, fmt.Errorf("login user: %w", err)
	}

	return &token, nil
}


func (a *AuthService) RegisterUser(ctx context.Context, username, email, password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generating password: %w", ErrorHashingPassword)
	}

	user := models.User{Username: username, Email: email, Password: hashedPass}
	id, err := a.Storage.CreateUser(context.Background(), user)
	if err != nil {
		return "", fmt.Errorf("registration user: %w", err)
	}

	user.Id = id
	token, err := my_jwt.NewToken(user)
	if err != nil {
		return "", fmt.Errorf("registration user: %w", err)
	}

	a.Kafka.PushEmailToQueue(ctx, "registration", email)

	return token, nil
}