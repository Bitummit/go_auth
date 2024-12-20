package grpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	auth "github.com/Bitummit/go_auth/internal/service"
	"github.com/Bitummit/go_auth/internal/storage/postgresql"
	my_jwt "github.com/Bitummit/go_auth/internal/jwt"
	auth_proto "github.com/Bitummit/go_auth/pkg/auth_proto_gen/proto"
	"github.com/Bitummit/go_auth/pkg/config"
	"github.com/Bitummit/go_auth/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type (
	AuthServer struct {
		Cfg *config.Config
		Log *slog.Logger
		Service Service
		auth_proto.UnimplementedAuthServer
	}

	Service interface {
		CheckTokenUser(token string) error
		LoginUser(cusername string, password string) (*string, error)
		RegisterUser(ctx context.Context, username string, email string, password string) (string, error)
	}
)


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
		if errors.Is(err, my_jwt.ErrorTokenDuration) || errors.Is(err, postgresql.ErrorNotFound){
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
		if errors.Is(err, postgresql.ErrorNotFound) || errors.Is(err, auth.ErrorHashingPassword){
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("%v", err))
	}

	response := auth_proto.LoginResponse{
		Token: *token,
	}
	return &response, nil
}

func (a *AuthServer) Register(ctx context.Context, req *auth_proto.RegistrationRequest) (*auth_proto.RegistrationResponse, error)  {
	token, err := a.Service.RegisterUser(ctx, req.GetUsername(), req.GetEmail(), req.GetPassword())
	if err != nil {
		a.Log.Error("error while register user:", logger.Err(err))
		if errors.Is(err, postgresql.ErrorUserExists) || errors.Is(err, auth.ErrorHashingPassword){
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("%v", err))
	}
	
	response := auth_proto.RegistrationResponse{
		Token: token,
	}
	return &response, nil
}

// https://www.koyeb.com/tutorials/build-a-grpc-api-using-go-and-grpc-gateway