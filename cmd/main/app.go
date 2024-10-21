package main

import (
	"context"
	"net"
	"os/signal"
	// "sync"

	"syscall"

	my_grpc "github.com/Bitummit/go_auth/internal/api/grpc"
	"github.com/Bitummit/go_auth/internal/service"
	"github.com/Bitummit/go_auth/internal/storage/postgresql"
	auth_proto "github.com/Bitummit/go_auth/pkg/auth_proto_gen/proto"
	"github.com/Bitummit/go_auth/pkg/config"
	"github.com/Bitummit/go_auth/pkg/logger"
	"google.golang.org/grpc"
)


func main() {
	ctx, stop  := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	// wg := &sync.WaitGroup{}
	// wg.Add(1)
	cfg := config.InitConfig()
	log := logger.NewLogger()
	log.Info("Initializing config success")

	log.Info("Connecting database ...")
	
	storage, err := postgresql.New(ctx)

	if err != nil {
		log.Error("Error connecting to DB", logger.Err(err))
		return
	}
	log.Info("Connecting database SUCCESS")
	service := service.New(storage, log)

	log.Info("starting server ...")
	// wg.Add(1)
	server := my_grpc.New(log, cfg, service)
	go startServer(ctx, server) 

	<-ctx.Done()
	storage.DB.Close()
	log.Info("Database stopped")
}


func startServer(ctx context.Context, server *my_grpc.AuthServer) {
	
	listener, err := net.Listen("tcp", server.Cfg.GrpcAddress)
	if err != nil {
		server.Log.Error("failed to listen", logger.Err(err))
	}
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	auth_proto.RegisterAuthServer(grpcServer, server)

	go func() {
		if err = grpcServer.Serve(listener); err != nil {
			server.Log.Error("error starting server", logger.Err(err))
		}
	}()
	<-ctx.Done()
	grpcServer.GracefulStop()
	server.Log.Info("Server stopped")
}
