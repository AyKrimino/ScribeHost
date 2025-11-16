package errors

import "fmt"

type InvalidTokenError struct {
	TokenType string
	Msg       string
}

func (e InvalidTokenError) Error() string {
	if e.TokenType != "" && e.Msg != "" {
		return fmt.Sprintf("%s: %s", e.TokenType, e.Msg)
	}
	if e.Msg != "" {
		return e.Msg
	}
	return "invalid token"
}

func IsInvalidTokenError(err error) bool {
	_, ok := err.(InvalidTokenError)
	return ok
}

func NewInvalidTokenError(tokenType, msg string) InvalidTokenError {
	return InvalidTokenError{
		TokenType: tokenType,
		Msg:       msg,
	}
}
