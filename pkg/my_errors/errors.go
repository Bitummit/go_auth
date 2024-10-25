package my_errors

import "errors"

var ErrorInvalidToken = errors.New("invalid token")
var ErrorTokenDuration = errors.New("invalid token duration")

var ErrorSigningToken = errors.New("token signing error")
