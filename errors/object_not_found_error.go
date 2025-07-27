package errors

import "fmt"

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
