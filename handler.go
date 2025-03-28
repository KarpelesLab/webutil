package webutil

import "net/http"

// handlerFunc is similar to http.HandlerFunc but also implements the error interface.
// This allows it to be used as both an HTTP handler and as an error return value
// for functions that might need to redirect or delegate handling.
type handlerFunc func(w http.ResponseWriter, req *http.Request)

// ServeHTTP implements the http.Handler interface by calling the function itself.
func (f handlerFunc) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	f(w, req)
}

// Error implements the error interface, returning a descriptive message.
// This allows the handlerFunc to be returned as an error.
func (f handlerFunc) Error() string {
	return "request must be handled by a different handler"
}

// ErrorHandler converts a standard http.Handler into an error type
// that still implements the http.Handler interface.
//
// This is useful for returning handlers as errors in functions that
// need to indicate that a request should be handled differently.
func ErrorHandler(handler http.Handler) error {
	return handlerFunc(handler.ServeHTTP)
}

// ErrorHandlerFunc converts a standard handler function into an error type
// that also implements the http.Handler interface.
//
// This allows function-based handlers to be returned as errors from
// functions that might need to delegate handling.
func ErrorHandlerFunc(handlerFn func(w http.ResponseWriter, req *http.Request)) error {
	return handlerFunc(handlerFn)
}
