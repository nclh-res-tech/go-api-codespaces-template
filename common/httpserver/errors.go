package httpserver

import (
	"net/http"
	"strings"

	coreerrors "{{MODULE_PATH}}/common/errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Error = coreerrors.Error

const (
	// ErrUnknown is returned when an unexpected error occurs.
	ErrUnknown = coreerrors.ErrUnknown
	// ErrInvalidRequest is returned when either the parameters or the request body is invalid.
	ErrInvalidRequest = coreerrors.ErrInvalidRequest
	// ErrValidation is returned when the parameters don't pass validation.
	ErrValidation = coreerrors.ErrValidation
	// ErrNotFound is returned when the requested resource is not found.
	ErrNotFound = coreerrors.ErrNotFound
	// ErrUnauthorized is returned when the request is not authorized.
	ErrUnauthorized = coreerrors.ErrUnauthorized
	// ErrSeperator is used to determine the boundaries of the errors in the hierarchy.
	ErrSeperator = coreerrors.ErrSeperator
)

func NewError(message string) error {
	return coreerrors.New(message)
}

func Is(err error, target error) bool {
	return coreerrors.Is(err, target)
}

func As(err error, target any) bool {
	return coreerrors.As(err, target)
}

func HandleResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func HandleError(c *gin.Context, err error) {
	zap.L().Error("error occurred in request", zap.Error(err))

	status := http.StatusInternalServerError
	switch {
	case coreerrors.Is(err, ErrInvalidRequest):
		fallthrough
	case coreerrors.Is(err, ErrValidation):
		status = http.StatusBadRequest
	case coreerrors.Is(err, ErrUnauthorized):
		status = http.StatusUnauthorized
	case coreerrors.Is(err, ErrNotFound):
		status = http.StatusNotFound
	case coreerrors.Is(err, ErrUnknown):
		status = http.StatusInternalServerError
	default:
		status = http.StatusInternalServerError
	}

	errJSON := struct {
		Error string `json:"error"`
	}{
		Error: strings.Split(err.Error(), ErrSeperator)[0], // TODO we may need to strip additional error information
	}

	c.JSON(status, errJSON)
}
