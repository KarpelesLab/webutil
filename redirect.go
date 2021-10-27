package webutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"
)

type Redirect struct {
	URL  *url.URL
	Code int
}

func SendRedirect(w http.ResponseWriter, url string, code int) {
	w.Header().Set("Location", url)
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code) // http.StatusFound

	fmt.Fprintf(w, "You are being redirected to <a href=\"%s\">%s</a>. If you see this message, please manually follow the link.", html.EscapeString(url), html.EscapeString(url))
	// try various stuff to cause the redirect to happen in case header failed to happen
	if js, err := json.Marshal(url); err == nil {
		fmt.Fprintf(w, "<script language=\"javascript\">window.location = %s;</script>", js)
	}
	fmt.Fprintf(w, "<meta http-equiv=\"Refresh\" content=\"0; url=%s\"/>", html.EscapeString(url))
	fmt.Fprintf(w, "<meta http-equiv=\"Location\" content=\"%s\"/>", html.EscapeString(url))
}

// code can be one of http.StatusMovedPermanently or http.StatusFound or
// any 3xx http status code
func RedirectErrorCode(u *url.URL, code int) error {
	// generate a redirect error
	n := &Redirect{URL: new(url.URL), Code: code}
	// copy url
	*n.URL = *u

	return n
}

func RedirectError(u *url.URL) error {
	// generate a redirect error
	n := &Redirect{URL: new(url.URL), Code: http.StatusFound}
	// copy url
	*n.URL = *u

	return n
}

func (e *Redirect) Error() string {
	return fmt.Sprintf("Redirect required to %s", e.URL)
}

func (e *Redirect) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	SendRedirect(w, e.URL.String(), e.Code)
}

func (e *Redirect) HTTPStatus() int {
	return e.Code
}

func IsRedirect(e error) http.Handler {
	var r *Redirect
	if errors.As(e, &r) {
		return r
	}
	return nil
}
