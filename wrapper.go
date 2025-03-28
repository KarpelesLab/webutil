package webutil

import "net/http"

// Handler extends the standard http.Handler interface by allowing
// ServeHTTP to return an error. This enables more flexible error handling
// and allows handlers to delegate requests via error returns.
type Handler interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request) error
}

// Wrapper adapts the Handler interface to the standard http.Handler interface.
// It automatically handles any errors returned by the wrapped Handler by
// converting them to HTTP responses.
type Wrapper struct {
	Child Handler // The wrapped Handler that may return errors
}

// WrapFunc is a function type that implements the Handler interface.
// It's similar to http.HandlerFunc but can return an error.
type WrapFunc func(w http.ResponseWriter, req *http.Request) error

// ServeHTTP implements the http.Handler interface by calling the wrapped Handler
// and handling any returned errors by converting them to appropriate HTTP responses.
func (wrapper *Wrapper) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := wrapper.Child.ServeHTTP(w, req)
	if err != nil {
		// Convert the error to an HTTP handler and serve the response
		ErrorToHTTPHandler(err).ServeHTTP(w, req)
	}
}

// ServeHTTP implements the http.Handler interface for WrapFunc, allowing
// it to be used directly as an http.Handler while still returning errors.
func (wf WrapFunc) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := wf(w, req)
	if err != nil {
		// Convert the error to an HTTP handler and serve the response
		ErrorToHTTPHandler(err).ServeHTTP(w, req)
	}
}

// Wrap converts a Handler to a standard http.Handler by wrapping it
// in a Wrapper that will handle any returned errors.
func Wrap(h Handler) http.Handler {
	return &Wrapper{Child: h}
}
