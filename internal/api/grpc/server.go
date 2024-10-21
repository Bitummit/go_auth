package grpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Bitummit/go_auth/pkg/my_errors"
	auth_proto "github.com/Bitummit/go_auth/pkg/auth_proto_gen/proto"
	"github.com/Bitummit/go_auth/pkg/config"
	"github.com/Bitummit/go_auth/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
type Service interface {
	CheckTokenUser(string) error
	LoginUser(string, string) (*string, error)
	RegisterUser(string, string) (*string, error)
}

type AuthServer struct {
	Cfg *config.Config
	Log *slog.Logger
	Service Service
	auth_proto.UnimplementedAuthServer
}

func New(log *slog.Logger, cfg *config.Config, service Service) *AuthServer {
	return &AuthServer{
		Cfg: cfg,
		Log: log,
		Service: service,
	}
}


func (a *AuthServer) CheckToken(ctx context.Context, req *auth_proto.CheckTokenRequest) (*auth_proto.EmptyResponse, error) {
	// Тут возвращать текст ошибки юзеру, а ниже можно подробнее
	if err := a.Service.CheckTokenUser(req.GetToken()); err != nil {
		a.Log.Error("error while login:", logger.Err(err))
		if errors.Is(err, my_errors.ErrorTokenDuration) || errors.Is(err, my_errors.ErrorNotFound){
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
			return nil, status.Error(codes.Internal, fmt.Sprintf("%v", err))
	}

	response := auth_proto.EmptyResponse{}
	return &response, nil
}


func (a *AuthServer) Login(ctx context.Context, req *auth_proto.LoginRequest) (*auth_proto.LoginResponse, error) {

	token, err := a.Service.LoginUser(req.GetUsername(), req.GetPassword())
	if err != nil {
		a.Log.Error("error while login:", logger.Err(err))
		if errors.Is(err, my_errors.ErrorNotFound) || errors.Is(err, my_errors.ErrorHashingPassword){
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("%v", err))
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
		if errors.Is(err, my_errors.ErrorUserExists) || errors.Is(err, my_errors.ErrorHashingPassword){
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("%v", err))
	}

	response := auth_proto.RegistrationResponse{
		Token: *token,
	}
	return &response, nil
}


// https://www.koyeb.com/tutorials/build-a-grpc-api-using-go-and-grpc-gateway