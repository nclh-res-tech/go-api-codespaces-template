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
