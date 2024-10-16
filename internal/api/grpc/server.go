package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/Bitummit/go_auth/internal/service"
	auth_proto "github.com/Bitummit/go_auth/pkg/auth_proto_gen/proto"
	"github.com/Bitummit/go_auth/pkg/config"
	"github.com/Bitummit/go_auth/pkg/logger"
	"google.golang.org/grpc"
)
type Service interface {
	CheckTokenUser(string) (bool, error)
	LoginUser(string, string) (*string, error)
	RegisterUser(string, string) (*string, error)
}

type AuthServer struct {
	Cfg *config.Config
	Log *slog.Logger
	Service Service
	auth_proto.UnimplementedAuthServer
}

func RunServer(log *slog.Logger, service Service, cfg *config.Config) error {
	
	server := AuthServer {
		Cfg: cfg,
		Log: log,
		Service: service,
	}
	server.Log.Info("starting server ...")
	listener, err := net.Listen("tcp", server.Cfg.GrpcAddress)
    if err != nil {
        server.Log.Error("failed to listen", logger.Err(err))
    }
	// ------- This to main ----------
    opts := []grpc.ServerOption{}
    grpcServer := grpc.NewServer(opts...)

	auth_proto.RegisterAuthServer(grpcServer, &server)
    if err = grpcServer.Serve(listener); err != nil {
		server.Log.Error("error starting server", logger.Err(err))
		return err
	}
	server.Log.Info("server stopped")
	return nil
	// ------- This to main ----------
}


func (a *AuthServer) CheckToken(ctx context.Context, req *auth_proto.CheckTokenRequest) (*auth_proto.EmptyResponse, error) {

	ok, err := a.Service.CheckTokenUser(req.GetToken())
	if err != nil || !ok {
		a.Log.Error("error while login:", logger.Err(err))
		return nil, fmt.Errorf("error in login: %v", err)
	}

	response := auth_proto.EmptyResponse{
	}

	return &response, nil
}


func (a *AuthServer) Login(ctx context.Context, req *auth_proto.LoginRequest) (*auth_proto.LoginResponse, error) {

	token, err := a.Service.LoginUser(req.GetUsername(), req.GetPassword())
	if err != nil {
		a.Log.Error("error while login:", logger.Err(err))
		return nil, err
	}

	response := auth_proto.LoginResponse{
		Token: *token,
	}
	return &response, nil
}


func (a *AuthServer) RegisterUser(lctx context.Context, req *auth_proto.RegistrationRequest) (*auth_proto.RegistrationResponse, error)  {

	token, err := a.Service.RegisterUser(req.GetUsername(), req.GetPassword())
	if err != nil {
		a.Log.Error("error while register user:", logger.Err(err))
		return nil, fmt.Errorf("error in login: %v", err)
	}

	response := auth_proto.RegistrationResponse{
		Token: *token,
	}
	return &response, nil
}


// https://www.koyeb.com/tutorials/build-a-grpc-api-using-go-and-grpc-gateway