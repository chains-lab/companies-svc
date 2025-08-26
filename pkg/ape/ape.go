package ape

import (
	"errors"

	"github.com/google/jsonapi"
)

type Error struct {
	// ID unique error identifier
	// in uppercase format like "ADMIN_CAN_NOT_DELETE_SELF"
	ID string

	// internal error which caused this error
	Cause error

	//Response error
	Response *jsonapi.ErrorObject
}

func (e *Error) Error() string {
	return e.ID
}

func (e *Error) Unwrap() error {
	if e.Cause != nil {
		return e.Cause
	}
	return nil
}

func (e *Error) Is(target error) bool {
	var be *Error
	if errors.As(target, &be) {
		return e.ID == be.ID
	}
	return false
}

func (e *Error) Raise(cause error, response *jsonapi.ErrorObject) error {
	return &Error{
		ID:       e.ID,
		Cause:    cause,
		Response: response,
	}
}

func (e *Error) JSONAPIError() *jsonapi.ErrorObject {
	return e.Response
}

func Declare(ID string) *Error {
	return &Error{
		ID: ID,
	}
}

func Unwrap(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return nil
}
