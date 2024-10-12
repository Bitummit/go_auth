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
}


func NewClient(log *slog.Logger, cfg *config.Config) (auth_v1.AuthClient, error) {

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
	defer conn.Close()

	client := auth_v1.NewAuthClient(conn)
	authClient.Client = client

	// request := &pb.Request{
	// 	Message: args[1],
	// }
	return client, nil


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

	return *token
}


// client := r;lgmerl
// client.Login()