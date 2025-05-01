// Package webutil provides utility functions and types for web applications.
package webutil

import (
	"net/http"
	"strings"
)

// SkipPrefix is an http.Handler that removes a prefix from the request URL path
// before passing the request to the underlying handler.
type SkipPrefix struct {
	// Prefix is the string to be removed from the beginning of the URL path
	Prefix string
	// Handler is the http.Handler that will serve the request after the prefix is removed
	Handler http.Handler
}

// AddPrefix is an http.Handler that adds a prefix to the request URL path
// before passing the request to the underlying handler.
type AddPrefix struct {
	// Prefix is the string to be added to the beginning of the URL path
	Prefix string
	// Handler is the http.Handler that will serve the request after the prefix is added
	Handler http.Handler
}

// ServeHTTP implements the http.Handler interface for SkipPrefix.
// It removes the specified prefix from the request URL path, RawPath, and RequestURI,
// sets a Sec-Access-Prefix header with the removed prefix,
// and then passes the modified request to the underlying handler.
func (h *SkipPrefix) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// trim prefix from request
	r.URL.Path = strings.TrimPrefix(r.URL.Path, h.Prefix)
	r.URL.RawPath = strings.TrimPrefix(r.URL.RawPath, h.Prefix)
	r.RequestURI = strings.TrimPrefix(r.RequestURI, h.Prefix)
	r.Header.Set("Sec-Access-Prefix", h.Prefix)
	// and serve
	h.Handler.ServeHTTP(w, r)
}

// pathJoin joins path segments, handling slashes appropriately.
// It ensures that there are no double slashes between segments.
// Unlike path.Join, it does not clean the path and preserves trailing slashes.
func pathJoin(p ...string) string {
	res := ""
	hasSlash := false

	for _, s := range p {
		if len(s) == 0 {
			continue
		}
		if s[0] == '/' && hasSlash {
			s = s[1:]
		}
		res = res + s
		hasSlash = strings.HasSuffix(res, "/")
	}

	return res
}

// ServeHTTP implements the http.Handler interface for AddPrefix.
// It adds the specified prefix to the request URL path, RawPath, and RequestURI,
// and then passes the modified request to the underlying handler.
func (h *AddPrefix) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// add prefix to request
	r.URL.Path = pathJoin(h.Prefix, r.URL.Path)
	r.URL.RawPath = pathJoin(h.Prefix, r.URL.RawPath)
	r.RequestURI = pathJoin(h.Prefix, r.RequestURI)
	// and serve
	h.Handler.ServeHTTP(w, r)
}
