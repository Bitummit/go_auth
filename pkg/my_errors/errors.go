package my_errors

import "errors"

var ErrorInvalidToken = errors.New("invalid token")
var ErrorTokenDuration = errors.New("invalid token duration")
var ErrorNotFound = errors.New("not found")
var ErrorWrongPassword = errors.New("wrong password")
var ErrorHashingPassword = errors.New("error while hashing password")
var ErrorSigningToken = errors.New("token signing error")
var ErrorUserExists = errors.New("user exists")
