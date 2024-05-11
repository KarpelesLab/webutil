package webutil

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
)

type serverError struct {
	e error
}

type HttpError int

// HttpErrorHandler returns a http.Handler for a given error code
//
// Deprecated: HttpError can be directly used as a handler
func HttpErrorHandler(code int) http.Handler {
	return HttpError(code)
}

// Error returns information about the error, including the corresponding text
func (e HttpError) Error() string {
	return fmt.Sprintf("HTTP error %d: %s", e, http.StatusText(int(e)))
}

// Unwrap will return a common error that (somewhat) matches the received error,
// making it easier to check if a given response matches a specific kind of error
func (e HttpError) Unwrap() error {
	switch e {
	case http.StatusBadRequest:
		return fs.ErrInvalid
	case http.StatusUnauthorized, http.StatusForbidden:
		return fs.ErrPermission
	case http.StatusNotFound:
		return fs.ErrNotExist
	default:
		return nil
	}
}

// ServeHTTP serves the error
func (e HttpError) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if e == http.StatusUnauthorized {
		w.Header().Set("WWW-Authenticate", "Basic realm=\"Website Access\"")
	}
	w.WriteHeader(int(e))
	fmt.Fprintf(w, "HTTP Error code %d: %s", int(e), http.StatusText(int(e)))
}

func (e *serverError) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	code := HTTPStatus(e.e)
	if code == 0 {
		code = http.StatusInternalServerError
	}

	w.WriteHeader(code)
	w.Write([]byte(e.e.Error()))
}

// HTTPStatus returns the value set in http error
func (e HttpError) HTTPStatus() int {
	return int(e)
}

func ErrorToHttpHandler(e error) http.Handler {
	if h, ok := e.(http.Handler); ok {
		return h
	}
	return &serverError{e}
}

type handlerError struct {
	h http.Handler
}

func HttpHandlerToError(h http.Handler) error {
	return &handlerError{h}
}

func (h *handlerError) Error() string {
	return "request is being forwarded"
}

func (h *handlerError) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.h.ServeHTTP(w, req)
}

func ServeError(w http.ResponseWriter, req *http.Request, err error) {
	var h http.Handler
	if errors.As(err, &h) {
		h.ServeHTTP(w, req)
		return
	}

	// fallback to instance of serverError
	(&serverError{err}).ServeHTTP(w, req)
}
