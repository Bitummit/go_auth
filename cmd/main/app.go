package main

import (
	"context"

	"github.com/Bitummit/go_auth/internal/grpc"
	"github.com/Bitummit/go_auth/internal/storage/postgresql"
	"github.com/Bitummit/go_auth/pkg/config"
	"github.com/Bitummit/go_auth/pkg/logger"

)


func main() {
	cfg := config.InitConfig()

	log := logger.NewLogger()
	log.Info("Initializing config success")

	log.Info("Connecting database ...")

	storage, err := postgresql.NewDBPool(context.TODO())
	if err != nil {
		log.Error("Error connecting to DB", logger.Err(err))
	}

	if err := grpc.StartGrpcServer(log, storage, cfg); err != nil {
		log.Error("Server error! Disconnecting ...")
	}

}