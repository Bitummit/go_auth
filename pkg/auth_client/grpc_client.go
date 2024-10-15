package authclient

import (
	"context"
	"log/slog"

	auth_v1 "github.com/Bitummit/go_auth/pkg/auth_v1/proto"
	"github.com/Bitummit/go_auth/pkg/config"
	"github.com/Bitummit/go_auth/pkg/logger"
	"google.golang.org/grpc"
)

type AuthClient struct {
	Client auth_v1.AuthClient
	Cfg *config.Config
	Log *slog.Logger
	Conn *grpc.ClientConn
}


func NewClient(log *slog.Logger, cfg *config.Config) (*AuthClient, error) {

	authClient := AuthClient {
		Cfg: cfg,
		Log: log,
	}

	opts := []grpc.DialOption{
	}
	
	conn, err := grpc.NewClient("127.0.0.1:5300", opts...)
	if err != nil {
		return nil, err
	}

	client := auth_v1.NewAuthClient(conn)
	authClient.Client = client
	authClient.Conn = conn

	return &authClient, nil
}


func (a *AuthClient) CheckToken(token string) auth_v1.Response {
	request := &auth_v1.Token {
		Token: token,
	}
	response, err := a.Client.CheckToken(context.Background(), request)
	if err != nil {
		a.Log.Error("fail to dial: %v", logger.Err(err))
	}
	// return response object or string?
	return *response
}


func (a *AuthClient) Login(username string, password string) auth_v1.Token {
	request := &auth_v1.BaseUserInformation {
		Username: username,
		Password: password,
	}
	token, err := a.Client.Login(context.Background(), request)
	if err != nil {
		a.Log.Error("fail to dial: %v", logger.Err(err))
	}
	// return Token object or string?
	return *token
}



func (a *AuthClient) Register(username string, password string) auth_v1.Token {
	request := &auth_v1.BaseUserInformation {
		Username: username,
		Password: password,
	}
	token, err := a.Client.Register(context.Background(), request)
	if err != nil {
		a.Log.Error("fail to dial: %v", logger.Err(err))
	}
	// return Token object or string?
	return *token
}
