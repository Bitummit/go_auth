package models


type User struct {
	Id int64
	Username string
	Email string
	Password []byte
}

type Token struct {
	Id int
	Access_token string
	Refresh_token string
}