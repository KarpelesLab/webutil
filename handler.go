package webutil

import "net/http"

// handlerFunc is similar to http.HandlerFunc but also complies with error type
// and can be used to return an error that means the request must be handled by
// something different (redirect, etc).
type handlerFunc func(w http.ResponseWriter, req *http.Request)

func (f handlerFunc) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	f(w, req)
}

func (f handlerFunc) Error() string {
	return "request must be handled by a different handler"
}

// ErrorHandler returns an error that also complies with http.Handler
// for a given http.Handler.
func ErrorHandler(f http.Handler) error {
	return handlerFunc(f.ServeHTTP)
}

// ErrorHandlerFunc returns an error that also complies with http.Handler
// for a given function.
func ErrorHandlerFunc(f func(w http.ResponseWriter, req *http.Request)) error {
	return handlerFunc(f)
}
