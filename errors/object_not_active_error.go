package errors

import "fmt"

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
