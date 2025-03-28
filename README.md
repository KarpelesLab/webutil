[![GoDoc](https://godoc.org/github.com/KarpelesLab/webutil?status.svg)](https://godoc.org/github.com/KarpelesLab/webutil)

# Web Utilities for Go

A collection of utilities for HTTP servers and web applications in Go.

## Features

- **HTTP Error Handling**: Type `HTTPError` implements both `error` and `http.Handler` interfaces
- **Resumable Downloads**: Auto-resuming HTTP downloads using Range headers
- **Data URI Parsing**: Parse and work with `data:` URI schemes
- **IP:Port Parsing**: Parse and validate IP addresses with ports
- **Error-to-Handler Conversion**: Convert between errors and HTTP handlers
- **PHP Query String Handling**: Parse and generate PHP-style query strings with array notation
- **HTTP Redirects**: Create and handle HTTP redirects as errors

## Usage Examples

```go
// Get with auto-resume
resp, err := webutil.Get("https://example.com/large-file.zip")
if err != nil {
    // Handle error
}
defer resp.Close()
// Reading from resp will auto-resume on network errors

// Return HTTP errors from handlers
if !authenticated {
    return webutil.StatusUnauthorized
}

// Parse data URIs
data, mime, err := webutil.ParseDataURI("data:text/plain;base64,SGVsbG8gV29ybGQ=")
```

## Documentation

See [GoDoc](https://godoc.org/github.com/KarpelesLab/webutil) for complete API documentation.

