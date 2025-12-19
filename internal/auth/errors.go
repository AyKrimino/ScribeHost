package auth

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

type InvalidOTPError struct {
	Msg string
}

func (e InvalidOTPError) Error() string {
	if e.Msg != "" {
		return fmt.Sprintf("invalid OTP: %s", e.Msg)
	}
	return "invlaid OTP"
}

func IsInvalidOTPError(err error) bool {
	_, ok := err.(InvalidOTPError)
	return ok
}

func NewInvalidOTPError(msg string) InvalidOTPError {
	return InvalidOTPError{
		Msg: msg,
	}
}

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

type ObjectAlreadyExistsError struct {
	ObjectType      string
	Identifier      string
	IdentifierValue string
}

func (e ObjectAlreadyExistsError) Error() string {
	if e.ObjectType != "" && e.Identifier != "" && e.IdentifierValue != "" {
		return fmt.Sprintf("%s already exists: %s %s", e.ObjectType, e.Identifier, e.IdentifierValue)
	}
	if e.ObjectType != "" {
		return fmt.Sprintf("%s already exists", e.ObjectType)
	}
	return "object already exists"
}

func IsObjectAlreadyExists(err error) bool {
	_, ok := err.(ObjectAlreadyExistsError)
	return ok
}

func NewObjectAlreadyExistsError(objectType, identifier, identifierValue string) ObjectAlreadyExistsError {
	return ObjectAlreadyExistsError{
		ObjectType:      objectType,
		Identifier:      identifier,
		IdentifierValue: identifierValue,
	}
}

type ObjectNotActiveError struct {
	ObjectType string
}

func (e ObjectNotActiveError) Error() string {
	if e.ObjectType != "" {
		return fmt.Sprintf("%s is not active", e.ObjectType)
	}
	return "object is not active"
}

func IsObjectNotActiveError(err error) bool {
	_, ok := err.(ObjectNotActiveError)
	return ok
}

func NewObjectNotActiveError(objectType string) ObjectNotActiveError {
	return ObjectNotActiveError{
		ObjectType: objectType,
	}
}

type ObjectNotFoundError struct {
	ObjectType      string
	Identifier      string
	IdentifierValue string
}

func (e ObjectNotFoundError) Error() string {
	if e.ObjectType != "" && e.Identifier != "" {
		return fmt.Sprintf("%s with this identifier %s %s not found", e.ObjectType, e.Identifier, e.IdentifierValue)
	}
	if e.ObjectType != "" {
		return fmt.Sprintf("%s not found", e.ObjectType)
	}
	return "object not found"
}

func IsObjectNotFoundError(err error) bool {
	_, ok := err.(ObjectNotFoundError)
	return ok
}

func NewObjectNotFoundError(objectType, identifier, identifierValue string) ObjectNotFoundError {
	return ObjectNotFoundError{
		ObjectType:      objectType,
		Identifier:      identifier,
		IdentifierValue: identifierValue,
	}
}
