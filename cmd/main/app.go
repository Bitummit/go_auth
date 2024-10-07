package main

import (
	"auth/internal/handlers"
	"auth/internal/storage/postgresql"
	"auth/pkg/config"
	"auth/pkg/logger"
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.URLFormat)

	router.Post("/register/", handlers.RegisterUser(log, storage))
	router.Post("/login/", handlers.LoginUser(log, storage))
	router.Post("/token/", handlers.CheckToken(log, storage))


	server := http.Server{
		Addr: cfg.HTTPServer.Address,
		Handler: router,
		ReadTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
	}
	log.Info("Starting server ", slog.String("addr", cfg.Address))

	if err := server.ListenAndServe(); err != nil {
		log.Error("Server error! Disconnecting ...")
	}

}