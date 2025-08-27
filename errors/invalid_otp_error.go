package errors

import "fmt"

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
