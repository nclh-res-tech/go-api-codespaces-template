package errors

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	// ErrUnknown is returned when an unexpected error occurs.
	ErrUnknown = Error("err_unknown: unknown error occurred")
	// ErrInvalidRequest is returned when either the parameters or the request body is invalid.
	ErrInvalidRequest = Error("err_invalid_request: invalid request received")
	// ErrValidation is returned when the parameters don't pass validation.
	ErrValidation = Error("err_validation: failed validation")
	// ErrNotFound is returned when the requested resource is not found.
	ErrNotFound = Error("err_not_found: not found")
	// ErrUnauthorized is returned when the request is not authorized.
	ErrUnauthorized = Error("err_unauthorized: unauthorized")
)

// ErrSeperator is used to determine the boundaries of the errors in the hierarchy.
const ErrSeperator = " -- "

// Error allows errors to be defined as const errors preventing modification
// and allowing them to be evaluated against wrapped errors.
type Error string

func (s Error) Error() string {
	return string(s)
}

// Is implements https://golang.org/pkg/errors/#Is allowing a Error
// to check it is the same even when wrapped. This implementation only
// checks the top most wrapped error.
func (s Error) Is(target error) bool {
	return s.Error() == target.Error() || strings.HasPrefix(target.Error(), s.Error()+ErrSeperator)
}

// As implements As(interface{}) bool which is used by errors.As
// allowing a Error to be set as the target if it matches the specified target type.
func (s Error) As(target interface{}) bool {
	v := reflect.ValueOf(target).Elem()
	if v.Type().Name() == "Error" && v.CanSet() {
		v.SetString(string(s))
		return true
	}
	return false
}

// Wrap allows errors to wrap an error returned from a 3rd party in a const service error.
func (s Error) Wrap(err error) error {
	return wrappedError{cause: err, msg: string(s)}
}

type wrappedError struct {
	cause error
	msg   string
}

func (w wrappedError) Error() string {
	if w.cause != nil {
		return fmt.Sprintf("%s%s%v", w.msg, ErrSeperator, w.cause)
	}
	return w.msg
}

func (w wrappedError) Is(target error) bool {
	return Error(w.msg).Is(target)
}

func (w wrappedError) As(target interface{}) bool {
	return Error(w.msg).As(target)
}

func (w wrappedError) Unwrap() error {
	return w.cause
}

// New just wraps errors.New.
func New(message string) error {
	//nolint:goerr113
	return errors.New(message)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}
