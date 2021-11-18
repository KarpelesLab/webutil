package webutil

import "net/http"

type Handler interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request) error
}

type Wrapper struct {
	Child Handler
}

func (wrapper *Wrapper) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := wrapper.Child.ServeHTTP(w, req)
	if err != nil {
		ErrorToHttpHandler(err).ServeHTTP(w, req)
	}
}
