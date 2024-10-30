package authclient

import (
	"context"
	"fmt"
	"log/slog"

	auth_gen "github.com/Bitummit/go_auth/pkg/auth_proto_gen/proto"
	"github.com/Bitummit/go_auth/pkg/config"
	"github.com/Bitummit/go_auth/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	Client auth_gen.AuthClient
	Cfg *config.Config
	Log *slog.Logger
}


func New(log *slog.Logger, cfg *config.Config) (*AuthClient, error) {
	authClient := AuthClient {
		Cfg: cfg,
		Log: log,
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient("127.0.0.1:5300", opts...)
	if err != nil {
		return nil, err
	}

	client := auth_gen.NewAuthClient(conn)
	authClient.Client = client

	return &authClient, nil
}

func (a *AuthClient) CheckToken(token string) (*auth_gen.EmptyResponse, error) {
	request := &auth_gen.CheckTokenRequest {
		Token: token,
	}
	response, err := a.Client.CheckToken(context.Background(), request)
	if err != nil {
		a.Log.Error("fail to dial: %v", logger.Err(err))
		return nil, fmt.Errorf("auth service error: %v", err)
	}
	return response, nil
}

func (a *AuthClient) Login(username string, password string) (*auth_gen.LoginResponse, error) {
	request := &auth_gen.LoginRequest {
		Username: username,
		Password: password,
	}
	response, err := a.Client.Login(context.Background(), request)
	if err != nil {
		a.Log.Error("fail to dial: %v", logger.Err(err))
		// return nil, fmt.Errorf("auth service error: %v", err)
		return nil, err
	}
	return response, nil
}

func (a *AuthClient) Register(username string, email string, password string) (*auth_gen.RegistrationResponse, error) {
	request := &auth_gen.RegistrationRequest {
		Username: username,
		Email: email,
		Password: password,
	}
	token, err := a.Client.Register(context.Background(), request)
	if err != nil {
		a.Log.Error("fail to dial: %v", logger.Err(err))
		return nil, fmt.Errorf("auth service error: %v", err)
	}
	return token, nil
}
