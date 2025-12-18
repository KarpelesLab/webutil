[![GoDoc](https://godoc.org/github.com/KarpelesLab/webutil?status.svg)](https://godoc.org/github.com/KarpelesLab/webutil)

# Web Utilities for Go

A collection of utilities for HTTP servers and web applications in Go.

## Installation

```bash
go get github.com/KarpelesLab/webutil
```

## Features

- **HTTP Error Handling**: Type `HTTPError` implements both `error` and `http.Handler` interfaces
- **Error-Returning Handlers**: Extended handler interface that allows returning errors
- **HTTP Redirects**: Represent redirects as errors for flexible control flow
- **Resumable Downloads**: Auto-resuming HTTP downloads using Range headers
- **Data URI Parsing**: Parse and decode RFC 2397 `data:` URI schemes
- **PHP Query String Parsing**: Parse and encode PHP-style query strings with array notation
- **URL Path Manipulation**: Add or remove prefixes from request paths
- **IP:Port Parsing**: Parse IP addresses with optional ports (IPv4 and IPv6)

## Usage Examples

### HTTP Error Handling

The `HTTPError` type can be used both as an error and as an HTTP handler:

```go
// Return HTTP errors from handlers
func handler(w http.ResponseWriter, req *http.Request) error {
    if !authenticated {
        return webutil.StatusUnauthorized
    }
    if resource == nil {
        return webutil.StatusNotFound
    }
    return nil
}

// Serve any error as HTTP response
if err != nil {
    webutil.ServeError(w, req, err)
    return
}

// Extract HTTP status from any error (including wrapped errors)
status := webutil.HTTPStatus(err)
```

### Error-Returning Handlers

Use `Handler` and `WrapFunc` for handlers that can return errors:

```go
// Function-based handler that returns errors
http.Handle("/", webutil.WrapFunc(func(w http.ResponseWriter, req *http.Request) error {
    data, err := fetchData()
    if err != nil {
        return webutil.StatusInternalServerError
    }
    json.NewEncoder(w).Encode(data)
    return nil
}))

// Interface-based handler
type myHandler struct{}

func (h *myHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) error {
    // Handle request, return error if needed
    return nil
}

http.Handle("/api", webutil.Wrap(&myHandler{}))
```

### HTTP Redirects as Errors

```go
func handler(w http.ResponseWriter, req *http.Request) error {
    if needsRedirect {
        target, _ := url.Parse("https://example.com/new-location")
        return webutil.RedirectError(target)           // 302 Found
        // or: return webutil.RedirectErrorCode(target, http.StatusMovedPermanently)
    }
    return nil
}
```

### Resumable Downloads

```go
// Get with auto-resume on connection failures
reader, err := webutil.Get("https://example.com/large-file.zip")
if err != nil {
    log.Fatal(err)
}
defer reader.Close()

// Reading will automatically resume if the connection drops
io.Copy(dst, reader)
```

### Data URI Parsing

```go
// Parse base64-encoded data URI
data, mimeType, err := webutil.ParseDataURI("data:text/plain;base64,SGVsbG8gV29ybGQ=")
// data: []byte("Hello World")
// mimeType: "text/plain"

// Parse URL-encoded data URI
data, mimeType, err := webutil.ParseDataURI("data:text/html,%3Ch1%3EHello%3C%2Fh1%3E")

// Get also supports data URIs
reader, err := webutil.Get("data:text/plain,Hello")
```

### PHP-Style Query String Parsing

```go
// Parse nested objects: "user[name]=John&user[age]=30"
result := webutil.ParsePhpQuery("user[name]=John&user[age]=30")
// result: map[string]any{"user": map[string]any{"name": "John", "age": "30"}}

// Parse arrays: "tags[]=go&tags[]=web"
result := webutil.ParsePhpQuery("tags[]=go&tags[]=web")
// result: map[string]any{"tags": []any{"go", "web"}}

// Convert from url.Values
values := url.Values{"items[]": {"a", "b"}}
result := webutil.ConvertPhpQuery(values)

// Encode back to query string
query := webutil.EncodePhpQuery(result)
```

### URL Path Manipulation

```go
// Strip prefix from requests before routing
http.Handle("/api/", &webutil.SkipPrefix{
    Prefix:  "/api",
    Handler: apiRouter,
})

// Add prefix to requests
http.Handle("/", &webutil.AddPrefix{
    Prefix:  "/v1",
    Handler: versionedHandler,
})
```

### IP:Port Parsing

```go
// Parse various IP:port formats
addr := webutil.ParseIPPort("127.0.0.1:8080")     // IPv4 with port
addr := webutil.ParseIPPort("192.168.1.1")        // IPv4 only
addr := webutil.ParseIPPort("[::1]:8080")         // IPv6 with port
addr := webutil.ParseIPPort("::1")                // IPv6 only
addr := webutil.ParseIPPort(":8080")              // Port only

if addr != nil {
    fmt.Printf("IP: %v, Port: %d\n", addr.IP, addr.Port)
}
```

## Pre-defined HTTP Status Errors

The package provides constants for all standard HTTP error status codes:

**4xx Client Errors:**
`StatusBadRequest`, `StatusUnauthorized`, `StatusPaymentRequired`, `StatusForbidden`, `StatusNotFound`, `StatusMethodNotAllowed`, `StatusNotAcceptable`, `StatusProxyAuthRequired`, `StatusRequestTimeout`, `StatusConflict`, `StatusGone`, `StatusLengthRequired`, `StatusPreconditionFailed`, `StatusRequestEntityTooLarge`, `StatusRequestURITooLong`, `StatusUnsupportedMediaType`, `StatusRequestedRangeNotSatisfiable`, `StatusExpectationFailed`, `StatusTeapot`, `StatusMisdirectedRequest`, `StatusUnprocessableEntity`, `StatusLocked`, `StatusFailedDependency`, `StatusTooEarly`, `StatusUpgradeRequired`, `StatusPreconditionRequired`, `StatusTooManyRequests`, `StatusRequestHeaderFieldsTooLarge`, `StatusUnavailableForLegalReasons`

**5xx Server Errors:**
`StatusInternalServerError`, `StatusNotImplemented`, `StatusBadGateway`, `StatusServiceUnavailable`, `StatusGatewayTimeout`, `StatusHTTPVersionNotSupported`, `StatusVariantAlsoNegotiates`, `StatusInsufficientStorage`, `StatusLoopDetected`, `StatusNotExtended`, `StatusNetworkAuthenticationRequired`

## Documentation

See [GoDoc](https://godoc.org/github.com/KarpelesLab/webutil) for complete API documentation.

## License

See LICENSE file.
