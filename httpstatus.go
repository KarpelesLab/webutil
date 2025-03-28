package webutil

import (
	"errors"
	"io/fs"
	"net/http"
)

// HTTPStatus extracts an HTTP status code from an error.
// 
// It can handle several types of errors:
// 1. Standard fs package errors (ErrNotExist, ErrPermission)
// 2. HTTPError type from this package
// 3. Any error that implements HTTPStatus() int
// 4. Wrapped errors (recursively checks unwrapped errors)
//
// Returns 0 if no status code can be determined from the error.
func HTTPStatus(err error) int {
	if err == nil {
		return 0
	}

	// Check standard fs errors
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return http.StatusNotFound
	case errors.Is(err, fs.ErrPermission):
		return http.StatusForbidden
	}

	// Try to extract status from various error types
	var httpErr HTTPError
	if errors.As(err, &httpErr) {
		return int(httpErr)
	}

	// Check for anything implementing HTTPStatus() method
	type statusGetter interface {
		HTTPStatus() int
	}
	var statusErr statusGetter
	if errors.As(err, &statusErr) {
		return statusErr.HTTPStatus()
	}

	// If we have a wrapped error, recursively unwrap and check
	var unwrapErr interface{ Unwrap() error }
	if errors.As(err, &unwrapErr) {
		if unwrapped := unwrapErr.Unwrap(); unwrapped != nil {
			return HTTPStatus(unwrapped)
		}
	}

	// No status found
	return 0
}
