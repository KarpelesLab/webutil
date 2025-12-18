// Package webutil provides utility functions and types for building HTTP servers
// and web applications in Go.
//
// # HTTP Error Handling
//
// The package provides a flexible error handling system that integrates with
// the standard net/http package. The [HTTPError] type represents HTTP status codes
// as errors and implements both the error and http.Handler interfaces:
//
//	// Return HTTP errors directly from handlers
//	if !authenticated {
//	    return webutil.StatusUnauthorized
//	}
//
//	// Serve any error as an HTTP response
//	webutil.ServeError(w, req, err)
//
// Pre-defined status constants are available for all standard 4xx and 5xx
// HTTP status codes (e.g., [StatusNotFound], [StatusInternalServerError]).
//
// The [HTTPStatus] function extracts HTTP status codes from any error, including
// standard library fs errors, wrapped errors, and custom error types.
//
// # Error-Returning Handlers
//
// The [Handler] interface extends http.Handler to allow ServeHTTP to return errors,
// enabling more flexible error handling patterns:
//
//	type myHandler struct{}
//
//	func (h *myHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) error {
//	    if err := doSomething(); err != nil {
//	        return webutil.StatusInternalServerError
//	    }
//	    return nil
//	}
//
//	// Wrap for use with standard http server
//	http.Handle("/", webutil.Wrap(&myHandler{}))
//
// The [WrapFunc] type provides a function-based alternative similar to http.HandlerFunc.
//
// # HTTP Redirects
//
// The [Redirect] type represents HTTP redirects as errors, allowing redirects to be
// returned from error-returning handlers:
//
//	func handler(w http.ResponseWriter, req *http.Request) error {
//	    targetURL, _ := url.Parse("https://example.com")
//	    return webutil.RedirectError(targetURL)
//	}
//
// Use [RedirectErrorCode] to specify a custom redirect status code.
//
// # Resumable HTTP Downloads
//
// The [Get] function provides automatic resume capability for HTTP downloads.
// If a connection is interrupted, it automatically resumes using HTTP Range headers:
//
//	reader, err := webutil.Get("https://example.com/large-file.zip")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer reader.Close()
//	io.Copy(dst, reader) // Automatically resumes on network interruption
//
// The function also supports data: URIs, decoding embedded content directly.
//
// # Data URI Parsing
//
// The [ParseDataURI] function parses RFC 2397 data URIs:
//
//	data, mimeType, err := webutil.ParseDataURI("data:text/plain;base64,SGVsbG8gV29ybGQ=")
//	// data: []byte("Hello World")
//	// mimeType: "text/plain"
//
// Both base64 and URL-encoded data URIs are supported.
//
// # PHP-Style Query String Parsing
//
// The package provides functions for parsing and encoding PHP-style query strings
// with array and object notation:
//
//	// Parse: "a[b]=c&a[d]=e" -> map[string]any{"a": map[string]any{"b": "c", "d": "e"}}
//	result := webutil.ParsePhpQuery("a[b]=c&a[d]=e")
//
//	// Parse: "items[]=foo&items[]=bar" -> map[string]any{"items": []any{"foo", "bar"}}
//	result := webutil.ParsePhpQuery("items[]=foo&items[]=bar")
//
//	// Encode back to query string
//	query := webutil.EncodePhpQuery(result)
//
// This is useful for interoperability with PHP applications or APIs that use
// this query string format.
//
// # URL Path Manipulation
//
// The [SkipPrefix] and [AddPrefix] handlers modify request URL paths before
// delegating to wrapped handlers:
//
//	// Remove "/api" prefix before routing
//	http.Handle("/api/", &webutil.SkipPrefix{
//	    Prefix:  "/api",
//	    Handler: apiRouter,
//	})
//
// # Network Utilities
//
// The [ParseIPPort] function parses IP addresses with optional port numbers,
// supporting both IPv4 and IPv6 formats:
//
//	addr := webutil.ParseIPPort("127.0.0.1:8080")    // IPv4 with port
//	addr := webutil.ParseIPPort("[::1]:8080")        // IPv6 with port
//	addr := webutil.ParseIPPort("192.168.1.1")       // IPv4 without port
package webutil
