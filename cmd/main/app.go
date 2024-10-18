package main

import (
	"context"
	"net"
	"os"

	my_grpc "github.com/Bitummit/go_auth/internal/api/grpc"
	"github.com/Bitummit/go_auth/internal/service"
	"github.com/Bitummit/go_auth/internal/storage/postgresql"
	auth_proto "github.com/Bitummit/go_auth/pkg/auth_proto_gen/proto"
	"github.com/Bitummit/go_auth/pkg/config"
	"github.com/Bitummit/go_auth/pkg/logger"
	"google.golang.org/grpc"
)


func main() {
	ctx, cancel  := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.InitConfig()
	log := logger.NewLogger()
	log.Info("Initializing config success")

	log.Info("Connecting database ...")
	storage, err := postgresql.New(ctx)
	if err != nil {
		log.Error("Error connecting to DB", logger.Err(err))
		os.Exit(1)
	}
	log.Info("Connecting database SUCCESS")
	service := service.New(storage)

	log.Info("starting server ...")
	server := my_grpc.New(log, cfg, service)
	listener, err := net.Listen("tcp", server.Cfg.GrpcAddress)
    if err != nil {
        server.Log.Error("failed to listen", logger.Err(err))
		
    }
    opts := []grpc.ServerOption{}
    grpcServer := grpc.NewServer(opts...)
	auth_proto.RegisterAuthServer(grpcServer, server)
	
    if err = grpcServer.Serve(listener); err != nil {
		server.Log.Error("error starting server", logger.Err(err))
		os.Exit(1)
	}

}
