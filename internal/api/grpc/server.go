package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/Bitummit/go_auth/internal/service"
	"github.com/Bitummit/go_auth/internal/storage"
	auth_v1 "github.com/Bitummit/go_auth/pkg/auth_v1/proto"
	"github.com/Bitummit/go_auth/pkg/config"
	"github.com/Bitummit/go_auth/pkg/logger"
	"google.golang.org/grpc"
)


type AuthServer struct {
	Cfg *config.Config
	Log *slog.Logger
	Storage storage.QueryFunctions
	auth_v1.UnimplementedAuthServer
}

func RunServer(log *slog.Logger, storage storage.QueryFunctions, cfg *config.Config) error {
	
	server := AuthServer {
		Cfg: cfg,
		Storage: storage,
		Log: log,
	}
	server.Log.Info("starting server ...")
	listener, err := net.Listen("tcp", server.Cfg.GrpcAddress)
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


func (a *AuthServer) CheckToken(ctx context.Context, req *auth_v1.CheckTokenRequest) (*auth_v1.EmptyResponse, error) {

	ok, err := service.CheckTokenUser(a.Storage, req.GetToken())
	if err != nil || !ok {
		a.Log.Error("error while login:", logger.Err(err))
		return nil, fmt.Errorf("error in login: %v", err)
	}

	response := auth_v1.EmptyResponse{
	}

	return &response, nil
}


func (a *AuthServer) Login(ctx context.Context, req *auth_v1.LoginRequest) (*auth_v1.LoginResponse, error) {

	token, err := service.LoginUser(a.Storage, req.GetUsername(), req.GetPassword())
	if err != nil {
		a.Log.Error("error while login:", logger.Err(err))
		return nil, err
	}

	response := auth_v1.LoginResponse{
		Token: *token,
	}
	return &response, nil
}


func (a *AuthServer) RegisterUser(lctx context.Context, req *auth_v1.RegistrationRequest) (*auth_v1.RegistrationResponse, error)  {

	token, err := service.RegisterUser(a.Storage, req.GetUsername(), req.GetPassword())
	if err != nil {
		a.Log.Error("error while register user:", logger.Err(err))
		return nil, fmt.Errorf("error in login: %v", err)
	}

	response := auth_v1.RegistrationResponse{
		Token: *token,
	}
	return &response, nil
}

// Pass logger
// only server?
// Project structure
// https://www.koyeb.com/tutorials/build-a-grpc-api-using-go-and-grpc-gateway