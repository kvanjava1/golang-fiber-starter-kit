package response

import (
	"errors"
	"net/http"
)

// Sentinel domain errors
var (
	ErrNotFound     = errors.New("resource not found")
	ErrConflict     = errors.New("resource conflict or duplicate entry")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrValidation   = errors.New("validation error")
	ErrBadRequest   = errors.New("bad request")
)

// MapError memetakan domain error ke HTTP status code dan pesan publik yang aman
func MapError(err error) (int, string) {
	if err == nil {
		return http.StatusOK, "OK"
	}

	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound, "Resource not found"
	case errors.Is(err, ErrConflict):
		return http.StatusConflict, "Resource already exists"
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized, "Unauthorized access"
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden, "Access forbidden"
	case errors.Is(err, ErrValidation), errors.Is(err, ErrBadRequest):
		return http.StatusBadRequest, "Validation failed"
	default:
		return http.StatusInternalServerError, "Internal server error"
	}
}
