package errors

import (
	"errors"
	"net/http"
)

const (
	// CodeUnknown is the default error code
	C0   = 0
	C400 = http.StatusBadRequest
	C403 = http.StatusForbidden
	C404 = http.StatusNotFound
	C500 = http.StatusInternalServerError
)

type Errors struct {
	// Code is the error code
	Code int `json:"code"`
	// Message is the error message
	Message string `json:"message"`
}

func (e *Errors) Error() string {
	return e.Message
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func New(message string, code ...int) *Errors {
	var (
		c = C0
		m = message
	)
	if len(code) > 0 {
		c = code[0]
	}
	if len(code) == 0 {
		c = C500
	}
	return &Errors{
		Code:    c,
		Message: m,
	}
}
