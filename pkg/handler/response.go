package handler

import auth_v1 "github.com/Bitummit/go_auth/pkg/auth_v1/proto"


type Response struct {
	Status string `json:"status"`
	Error string `json:"error,omitempty"`
}


func OK() *auth_v1.Response {
	return &auth_v1.Response{
		Status: "OK",
	}
}

func Error() *auth_v1.Response {
	return &auth_v1.Response{
		Status: "Error",
	}
}