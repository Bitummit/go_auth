package service

import (
	"context"
	"fmt"

	"github.com/Bitummit/go_auth/internal/storage"
	"github.com/Bitummit/go_auth/internal/utils"
	"golang.org/x/crypto/bcrypt"
)


func CheckTokenUserService(queryTool storage.QueryFunctions, token string) (bool, error) {
	
	user, err := utils.ParseToken(token)
	if err != nil {
		return false, fmt.Errorf("wrong token %v", err)
	}
	
	_, err = queryTool.GetUser(context.Background(), user.Username)
	if err != nil {
		return false, fmt.Errorf("no such user %v", err)
	}

	return true, nil
}


func LoginUserService(queryTool storage.QueryFunctions, username string, password string) (*string, error) {
	user, err := queryTool.GetUser(context.Background(), username)
	if err != nil {
		// return nil, fmt.Errorf("error while fething data %v", err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password)); if err != nil {
		// return nil, fmt.Errorf("wrong password %v", err)
		return nil, err
	}

	token, err := utils.NewToken(*user)
	if err != nil {
		// return nil, fmt.Errorf("error while generating token %v", err)
		return nil, err
	}
	return &token, nil
}


func RegisterUserService(queryTool storage.QueryFunctions, username string, password string) (*string, error) {
	
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error while hashing password %v", err)
	}
	user := storage.User{Username: username, Password: hashedPass}
	id, err := queryTool.CreateUser(context.Background(), user)
	if err != nil {
		return nil, fmt.Errorf("error while inserting user %v", err)
	}
	user.Id = id
	token, err := utils.NewToken(user)
	if err != nil {
		return nil, fmt.Errorf("error while generating token %v", err)
	}

	return &token, nil
}