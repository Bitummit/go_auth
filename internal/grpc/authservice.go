package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/Bitummit/go_auth/internal/storage"
	"github.com/Bitummit/go_auth/internal/utils"
	auth_v1 "github.com/Bitummit/go_auth/pkg/auth_v1/proto"
	"github.com/Bitummit/go_auth/pkg/config"
	"github.com/Bitummit/go_auth/pkg/handler"
	"github.com/Bitummit/go_auth/pkg/logger"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)


type AuthServer struct {
	Cfg *config.Config
	Log *slog.Logger
	Storage storage.QueryFunctions
	auth_v1.UnimplementedAuthServer
}

func StartGrpcServer(log *slog.Logger, storage storage.QueryFunctions, cfg *config.Config) error {
	
	server := AuthServer {
		Cfg: cfg,
		Storage: storage,
		Log: log,
	}
	server.Log.Info("starting server ...")
	listener, err := net.Listen("tcp", ":5300")
    if err != nil {
        server.Log.Error("failed to listen", logger.Err(err))
    }

    opts := []grpc.ServerOption{}
    grpcServer := grpc.NewServer(opts...)

	auth_v1.RegisterAuthServer(grpcServer, &server)
    if err = grpcServer.Serve(listener); err != nil {
		server.Log.Error("error starting server", logger.Err(err))
		return err
	}
	server.Log.Info("server stopped")
	return nil
}

func (a *AuthServer) Login(ctx context.Context, req *auth_v1.BaseUserInformation) (*auth_v1.Token, error) {
	user, err := a.Storage.GetUser(context.Background(), req.GetUsername())
	if err != nil {
		a.Log.Error("Error inserting user", logger.Err(err))
		response := auth_v1.Token{
			Response: handler.Error(),
			Token: "",
		}
		return &response, fmt.Errorf("Error while fething data")
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password)); if err != nil {
		a.Log.Error("Error on hashing password")
		response := auth_v1.Token{
			Response: handler.Error(),
			Token: "",
		}
		return &response, fmt.Errorf("Error while hashing password")
	}

	token, err := utils.NewToken(*user)
	if err != nil {
		a.Log.Error("Error creating token", logger.Err(err))
		response := auth_v1.Token{
			Response: handler.Error(),
			Token: "",
		}
		return &response, fmt.Errorf("Error while generating token")
	}
	response := auth_v1.Token{
		Response: handler.OK(),
		Token: token,
	}
	return &response, nil
}

// Pass logger
// only server?
// Project structure
// https://www.koyeb.com/tutorials/build-a-grpc-api-using-go-and-grpc-gateway