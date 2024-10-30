package run

import (
	"context"
	"fmt"
	"net"
	"os/signal"
	"sync"

	"syscall"

	my_grpc "github.com/Bitummit/go_auth/internal/api/grpc"
	my_kafka "github.com/Bitummit/go_auth/internal/api/kafka"
	auth "github.com/Bitummit/go_auth/internal/service"
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

	cfg := config.InitConfig()
	log := logger.NewLogger()
	log.Info("Initializing config success")

	log.Info("Connecting database ...")
	storage, err := postgresql.New(ctx)
	if err != nil {
		log.Error("Error connecting to DB", logger.Err(err))
		return
	}
	log.Info("Connecting database success")
	
	kafkaServer, err := startKafka(ctx)
	if err != nil {
		log.Error("starting kafka", logger.Err(err))
	}

	wg.Add(1)
	service := auth.New(storage, log, kafkaServer)
	log.Info("Starting server ...")
	server := my_grpc.New(log, cfg, service)
	go startServer(ctx, wg, server) 

	<-ctx.Done()
	kafkaServer.Conn.Close()
	kafkaServer.Writer.Close()
	log.Info("Kafka stopped")
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

func startKafka(ctx context.Context) (*my_kafka.Kafka, error){
	kafkaServer, err := my_kafka.New(ctx, "localhost:9092", "emails", 0, []string{"localshost:9092"})
	if err != nil {
		return nil ,fmt.Errorf("connecting to kafka: %v", err)
	}
	kafkaServer.InitProducer()

	return kafkaServer, nil
}