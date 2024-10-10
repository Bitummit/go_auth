package handlers

// import (
// 	"github.com/Bitummit/go_auth/internal/storage"
// 	"github.com/Bitummit/go_auth/internal/utils"
// 	"github.com/Bitummit/go_auth/pkg/handler"
// 	"github.com/Bitummit/go_auth/pkg/logger"
// 	"context"
// 	"log/slog"
// 	"net/http"

// 	"github.com/go-chi/render"
// 	"golang.org/x/crypto/bcrypt"
// )
// type RegisterUserResponse struct {
// 	Response handler.Response
// 	Token string `json:"token,omitempty"`
// }

// type RegisterUserRequest struct{
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }

// func RegisterUser(log *slog.Logger, queryTool storage.QueryFunctions) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var req RegisterUserRequest

// 		err := render.DecodeJSON(r.Body, &req)
// 		if err != nil {
// 			log.Error("Cannot decode body", logger.Err(err))
// 			w.WriteHeader(http.StatusBadRequest)
// 			render.JSON(w, r, handler.Error("internal error"))
// 			return
// 		}

// 		hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
// 		if err != nil {
// 			log.Error("Error on hashing password")
// 			w.WriteHeader(http.StatusBadRequest)
// 			render.JSON(w, r, handler.Error("internal error"))
// 			return
// 		}
// 		user := storage.User{Username: req.Username, Password: hashedPass}
// 		id, err := queryTool.CreateUser(context.Background(), user)
// 		if err != nil {
// 			log.Error("Error inserting user", logger.Err(err))
// 			w.WriteHeader(http.StatusInternalServerError)
// 			render.JSON(w, r, handler.Error("internal error"))
// 			return
// 		}
// 		user.Id = id
// 		token, err := utils.NewToken(user)
// 		if err != nil {
// 			log.Error("Error creating token", logger.Err(err))
// 			w.WriteHeader(http.StatusInternalServerError)
// 			render.JSON(w, r, handler.Error("internal error"))
// 			return
// 		}
// 		w.WriteHeader(http.StatusOK)
// 		render.JSON(w, r, RegisterUserResponse{Response: handler.OK(), Token: token})
// 	}
// }

// type LoginUserResponse struct {
// 	Response handler.Response
// 	Token string `json:"token,omitempty"`
// }

// type LoginUserRequest struct{
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }


// func LoginUser(log *slog.Logger, queryTool storage.QueryFunctions) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var req LoginUserRequest

// 		err := render.DecodeJSON(r.Body, &req)
// 		if err != nil {
// 			log.Error("Cannot decode body", logger.Err(err))
// 			w.WriteHeader(http.StatusBadRequest)
// 			render.JSON(w, r, handler.Error("internal error"))
// 			return
// 		}

// 		user, err := queryTool.GetUser(context.Background(), req.Username)
// 		if err != nil {
// 			log.Error("Error inserting user", logger.Err(err))
// 			w.WriteHeader(http.StatusBadRequest)
// 			render.JSON(w, r, handler.Error("Login or password incorrect"))
// 			return
// 		}

// 		err = bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password)); if err != nil {
// 			log.Error("Error on hashing password")
// 			w.WriteHeader(http.StatusBadRequest)
// 			render.JSON(w, r, handler.Error("Login or password incorrect"))
// 			return
// 		}

// 		token, err := utils.NewToken(*user)
// 		if err != nil {
// 			log.Error("Error creating token", logger.Err(err))
// 			w.WriteHeader(http.StatusInternalServerError)
// 			render.JSON(w, r, handler.Error("internal error"))
// 			return
// 		}
// 		w.WriteHeader(http.StatusOK)
// 		render.JSON(w, r, RegisterUserResponse{Response: handler.OK(), Token: token})
// 	}
// }

// type CheckTokenResponse struct {
// 	Response handler.Response
// }

// type CheckTokenRequest struct{
// 	Token string `json:"token"`
// }

// func CheckToken(log *slog.Logger, queryTool storage.QueryFunctions) http.HandlerFunc{
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var req CheckTokenRequest

// 		err := render.DecodeJSON(r.Body, &req)
// 		if err != nil {
// 			log.Error("Cannot decode body", logger.Err(err))
// 			w.WriteHeader(http.StatusBadRequest)
// 			render.JSON(w, r, handler.Error("internal error"))
// 			return
// 		}
// 		user, err := utils.ParseToken(req.Token)
// 		if err != nil {
// 			log.Error("Error token", logger.Err(err))
// 			w.WriteHeader(http.StatusBadRequest)
// 			render.JSON(w, r, handler.Error("invalid token"))
// 			return
// 		}
// 		_, err = queryTool.GetUser(context.Background(), user.Username)
// 		if err != nil {
// 			log.Error("No such user", logger.Err(err))
// 			w.WriteHeader(http.StatusBadRequest)
// 			render.JSON(w, r, handler.Error("No such user"))
// 			return
// 		}

// 		w.WriteHeader(http.StatusOK)
// 		render.JSON(w, r, handler.OK())

// 	}
// }