package errors

import "fmt"

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
