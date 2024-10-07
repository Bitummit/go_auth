package storage

import "context"

type User struct {
	Id int64
	Username string
	Password []byte
}

type Token struct {
	Id int
	Access_token string
	Refresh_token string
}


type QueryFunctions interface {
	CreateUser(context.Context, User) (int64, error)
	GetUser(context.Context, string) (*User, error)
}