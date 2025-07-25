package errors

import "fmt"

type ObjectAlreadyExistsError struct {
	ObjectType string
	Identifier string
}

func (e ObjectAlreadyExistsError) Error() string {
	if e.ObjectType != "" && e.Identifier != "" {
		return fmt.Sprintf("%s already exists: %s", e.ObjectType, e.Identifier)
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

func NewObjectAlreadyExistsError(objectType, identifier string) ObjectAlreadyExistsError {
	return ObjectAlreadyExistsError{
		ObjectType: objectType,
		Identifier: identifier,
	}
}
