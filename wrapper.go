package webutil

import "net/http"

type Handler interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request) error
}

type Wrapper struct {
	Child Handler
}

type WrapFunc func(w http.ResponseWriter, req *http.Request) error

func (wrapper *Wrapper) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := wrapper.Child.ServeHTTP(w, req)
	if err != nil {
		ErrorToHttpHandler(err).ServeHTTP(w, req)
	}
}

func (wf WrapFunc) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := wf(w, req)
	if err != nil {
		ErrorToHttpHandler(err).ServeHTTP(w, req)
	}
}

func Wrap(h Handler) http.Handler {
	return &Wrapper{h}
}
