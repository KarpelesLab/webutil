package webutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"
)

// Redirect represents an error that requires an HTTP redirect.
// It implements both error and http.Handler interfaces.
type Redirect struct {
	URL  *url.URL // Target URL
	Code int      // HTTP status code (should be 3xx)
}

// SendRedirect is a robust HTTP redirect implementation that uses
// multiple fallback techniques to ensure successful redirection.
//
// Deprecated: Use http.Redirect() for standard redirects. This function
// is only needed in edge cases where user agents may not support standard
// redirects.
func SendRedirect(w http.ResponseWriter, target string, code int) {
	// Set standard redirect headers
	w.Header().Set("Location", target)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	// Provide fallback content with multiple redirect techniques
	escapedTarget := html.EscapeString(target)

	// Link for manual navigation
	_, _ = fmt.Fprintf(w, "You are being redirected to <a href=\"%s\">%s</a>. "+
		"If you see this message, please manually follow the link.",
		escapedTarget, escapedTarget)

	// JavaScript redirect
	if js, err := json.Marshal(target); err == nil {
		_, _ = fmt.Fprintf(w, "<script>window.location = %s;</script>", js)
	}

	// Meta refresh redirect
	_, _ = fmt.Fprintf(w, "<meta http-equiv=\"Refresh\" content=\"0; url=%s\"/>", escapedTarget)
	_, _ = fmt.Fprintf(w, "<meta http-equiv=\"Location\" content=\"%s\"/>", escapedTarget)
}

// RedirectErrorCode creates a redirect error with a specific HTTP status code.
//
// Parameters:
//   - u: Target URL for the redirection
//   - code: HTTP status code (should be a 3xx status code)
//
// Returns a Redirect object implementing both error and http.Handler interfaces.
func RedirectErrorCode(u *url.URL, code int) error {
	// Create a new redirect error with a deep copy of the URL
	redirect := &Redirect{
		URL:  &url.URL{},
		Code: code,
	}

	*redirect.URL = *u // Deep copy the URL
	return redirect
}

// RedirectError creates a redirect error with the default 302 Found status.
//
// Parameter:
//   - u: Target URL for the redirection
//
// Returns a Redirect object implementing both error and http.Handler interfaces.
func RedirectError(u *url.URL) error {
	return RedirectErrorCode(u, http.StatusFound)
}

// Error implements the error interface for Redirect.
func (r *Redirect) Error() string {
	return fmt.Sprintf("Redirect required to %s", r.URL)
}

// ServeHTTP implements the http.Handler interface for Redirect.
// When called, it performs an HTTP redirect to the target URL.
func (r *Redirect) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, r.URL.String(), r.Code)
}

// HTTPStatus returns the HTTP status code for this redirect.
func (r *Redirect) HTTPStatus() int {
	return r.Code
}

// IsRedirect checks if an error is a Redirect error.
// If it is, returns the error as an http.Handler for convenient use.
// Otherwise, returns nil.
func IsRedirect(err error) http.Handler {
	var redirect *Redirect
	if errors.As(err, &redirect) {
		return redirect
	}
	return nil
}
