package webutil

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

type serverError struct {
	e error
}

type httpError int

func HttpError(code int) error {
	return httpError(code)
}

func HttpErrorHandler(code int) http.Handler {
	return httpError(code)
}

func (e httpError) Error() string {
	return fmt.Sprintf("HTTP error %d", e)
}

func (e httpError) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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

	switch e.e {
	case os.ErrNotExist:
		http.NotFound(w, req)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Server error: %s", e.e)
	}
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
