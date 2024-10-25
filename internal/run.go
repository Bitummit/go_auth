package run

import (
	"context"
	"net"
	"os/signal"
	"sync"

	"syscall"

	my_grpc "github.com/Bitummit/go_auth/internal/api/grpc"
	"github.com/Bitummit/go_auth/internal/service"
	"github.com/Bitummit/go_auth/internal/storage/postgresql"
	auth_proto "github.com/Bitummit/go_auth/pkg/auth_proto_gen/proto"
	"github.com/Bitummit/go_auth/pkg/config"
	"github.com/Bitummit/go_auth/pkg/logger"
	"google.golang.org/grpc"
)


func Run() {
	ctx, stop  := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	wg := &sync.WaitGroup{}
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
	wg.Add(1)
	server := my_grpc.New(log, cfg, service)
	go startServer(ctx, wg, server) 

	<-ctx.Done()
	wg.Wait()
	storage.DB.Close()
	log.Info("Database stopped")
}



func startServer(ctx context.Context, wg *sync.WaitGroup, server *my_grpc.AuthServer) {	
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
	defer wg.Done()
	grpcServer.GracefulStop()
	server.Log.Info("Server stopped")
}


// go func() {
// 	defer wg.Done()
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			grpcServer.GracefulStop()
// 		default:
// 			if err = grpcServer.Serve(listener); err != nil {
// 				server.Log.Error("error starting server", logger.Err(err))
// 			}
// 		}
// 	}
// }()