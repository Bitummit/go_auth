package storage

import (
	"context"

	"github.com/Bitummit/go_auth/internal/models"
)


type QueryFunctions interface {
	CreateUser(context.Context, models.User) (int64, error)
	GetUser(context.Context, string) (*models.User, error)
}