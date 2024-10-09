package grpc

import (
	"auth/internal/storage"
	auth_v1 "auth/pkg/auth_v1/proto"
	"context"
	"errors"
)


type AuthService struct {
	Storage *storage.QueryFunctions
	auth_v1.UnimplementedAuthServer
}

func NewAuthService(storage *storage.QueryFunctions) AuthService {
	return AuthService{Storage: storage}
}

func (a *AuthService) Login(ctx context.Context, req *auth_v1.BaseUserInformation) (*auth_v1.Token, error) {
	return nil, errors.New("Test")
}

// Pass logger
// only server?
// Project structure
// https://www.koyeb.com/tutorials/build-a-grpc-api-using-go-and-grpc-gateway