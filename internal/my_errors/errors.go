package my_errors

type InvalidTokenError struct {}

type SigningTokenError struct {}

type TokenDurationError struct {}

func (e *InvalidTokenError) Error() string{
	return "invalid token"
}

func (e *SigningTokenError) Error() string{
	return "token signing error"
}

func (e *TokenDurationError) Error() string{
	return "invalid token duration"
}