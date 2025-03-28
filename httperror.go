package webutil

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
)

// serverError wraps a generic error to be served via HTTP
type serverError struct {
	e error
}

// HTTPError represents an HTTP error code as an error type.
// It also implements the http.Handler interface to serve the error.
type HTTPError int

// HttpError is an alias for HTTPError to maintain backward compatibility.
// Deprecated: Use HTTPError instead.
type HttpError = HTTPError

// HTTPErrorHandler returns a http.Handler for a given error code.
//
// Deprecated: HTTPError can be directly used as a handler.
func HTTPErrorHandler(code int) http.Handler {
	return HTTPError(code)
}

// HttpErrorHandler is an alias for HTTPErrorHandler to maintain backward compatibility.
//
// Deprecated: Use HTTPErrorHandler instead.
func HttpErrorHandler(code int) http.Handler {
	return HTTPErrorHandler(code)
}

// Error returns a formatted string with the HTTP error code and text.
func (e HTTPError) Error() string {
	return fmt.Sprintf("HTTP error %d: %s", e, http.StatusText(int(e)))
}

// Unwrap maps HTTP errors to standard filesystem errors when appropriate,
// making it easier to check if a response matches a specific kind of error.
func (e HTTPError) Unwrap() error {
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

// ServeHTTP implements the http.Handler interface to serve the HTTP error.
func (e HTTPError) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// Add WWW-Authenticate header for 401 Unauthorized
	if e == http.StatusUnauthorized {
		w.Header().Set("WWW-Authenticate", "Basic realm=\"Website Access\"")
	}

	w.WriteHeader(int(e))
	_, _ = fmt.Fprintf(w, "HTTP Error code %d: %s", int(e), http.StatusText(int(e)))
}

// ServeHTTP implements http.Handler for serverError to serve an error via HTTP
func (e *serverError) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// Get HTTP status code from error or default to 500
	code := HTTPStatus(e.e)
	if code == 0 {
		code = http.StatusInternalServerError
	}

	w.WriteHeader(code)
	_, _ = io.WriteString(w, e.e.Error())
}

// HTTPStatus returns the numeric HTTP status code represented by this error.
func (e HTTPError) HTTPStatus() int {
	return int(e)
}

// ErrorToHTTPHandler converts an error to an http.Handler.
// If the error already implements http.Handler, it's returned directly.
// Otherwise, it's wrapped in a serverError.
func ErrorToHTTPHandler(e error) http.Handler {
	if h, ok := e.(http.Handler); ok {
		return h
	}
	return &serverError{e}
}

// ErrorToHttpHandler is an alias for ErrorToHTTPHandler to maintain backward compatibility.
//
// Deprecated: Use ErrorToHTTPHandler instead.
func ErrorToHttpHandler(e error) http.Handler {
	return ErrorToHTTPHandler(e)
}

// handlerError wraps an http.Handler as an error
type handlerError struct {
	h http.Handler
}

// HTTPHandlerToError converts an http.Handler to an error.
func HTTPHandlerToError(h http.Handler) error {
	return &handlerError{h}
}

// HttpHandlerToError is an alias for HTTPHandlerToError to maintain backward compatibility.
//
// Deprecated: Use HTTPHandlerToError instead.
func HttpHandlerToError(h http.Handler) error {
	return HTTPHandlerToError(h)
}

// Error implements the error interface for handlerError.
func (h *handlerError) Error() string {
	return "request is being forwarded"
}

// ServeHTTP delegates to the wrapped handler.
func (h *handlerError) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.h.ServeHTTP(w, req)
}

// ServeError serves an error via HTTP.
// If the error can be type-asserted to an http.Handler, it uses that.
// Otherwise, it wraps the error in a serverError and serves that.
func ServeError(w http.ResponseWriter, req *http.Request, err error) {
	var h http.Handler
	if errors.As(err, &h) {
		h.ServeHTTP(w, req)
		return
	}

	// Fallback to instance of serverError
	(&serverError{err}).ServeHTTP(w, req)
}
