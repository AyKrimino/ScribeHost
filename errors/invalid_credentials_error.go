package errors

import "fmt"

type InvalidCredentialsError struct {
	Identifier string
}

func (e InvalidCredentialsError) Error() string {
	if e.Identifier != "" {
		return fmt.Sprintf("invalid credentials: %s is invalid", e.Identifier)
	}
	return "invalid credentials"
}

func IsInvalidCredentialsError(err error) bool {
	_, ok := err.(InvalidCredentialsError)
	return ok
}

func NewInvalidCredentialsError(identifier string) InvalidCredentialsError {
	return InvalidCredentialsError{
		Identifier: identifier,
	}
}
