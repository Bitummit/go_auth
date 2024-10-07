package handlers

import (
	"auth/internal/storage"
	"auth/internal/utils"
	"auth/pkg/handler"
	"auth/pkg/logger"
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)
type RegisterUserResponse struct {
	Response handler.Response
	Token string `json:"token,omitempty"`
}

type RegisterUserRequest struct{
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterUser(log *slog.Logger, queryTool storage.QueryFunctions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterUserRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Cannot decode body", logger.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, handler.Error("internal error"))
			return
		}

		hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("Error on hashing password")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, handler.Error("internal error"))
			return
		}
		user := storage.User{Username: req.Username, Password: hashedPass}
		id, err := queryTool.CreateUser(context.Background(), user)
		if err != nil {
			log.Error("Error inserting user", logger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, handler.Error("internal error"))
			return
		}
		user.Id = id
		token, err := utils.NewToken(user)
		if err != nil {
			log.Error("Error creating token", logger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, handler.Error("internal error"))
			return
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, RegisterUserResponse{Response: handler.OK(), Token: token})
	}
}