package main

import (
	"context"

	"github.com/Bitummit/go_auth/internal/api/grpc"
	"github.com/Bitummit/go_auth/internal/storage/postgresql"
	"github.com/Bitummit/go_auth/pkg/config"
	"github.com/Bitummit/go_auth/pkg/logger"

)


func main() {
	cfg := config.InitConfig()

	log := logger.NewLogger()
	log.Info("Initializing config success")

	log.Info("Connecting database ...")

	_, err := postgresql.NewDBPool(context.TODO())
	if err != nil {
		log.Error("Error connecting to DB", logger.Err(err))
	}
	// service := new
	// server := new
	if err := grpc.RunServer(log, cfg); err != nil {
		log.Error("Server error! Disconnecting ...")
	}

}