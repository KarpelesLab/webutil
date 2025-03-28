package webutil

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// ErrNotDataURI is returned when the input string is not a valid data URI
var ErrNotDataURI = errors.New("not a data URI")

// ErrNoEncodedValue is returned when the data URI doesn't contain a value part
var ErrNoEncodedValue = errors.New("could not locate encoded value")

// ParseDataURI parses a given data: URI and returns its binary data and MIME type.
//
// The format of a data URI is:
//
//	data:[<media type>][;base64],<data>
//
// Example: "data:text/plain;base64,SGVsbG8gV29ybGQ="
func ParseDataURI(uri string) ([]byte, string, error) {
	// Validate URI prefix
	if !strings.HasPrefix(uri, "data:") {
		return nil, "", ErrNotDataURI
	}

	// Remove the "data:" prefix
	uri = uri[5:]

	// Skip extra leading slashes (some implementations include them)
	uri = strings.TrimLeft(uri, "/")

	// Find the comma that separates metadata from the data
	p := strings.IndexByte(uri, ',')
	if p == -1 {
		return nil, "", ErrNoEncodedValue
	}

	// Split the metadata part into options (MIME type and encoding info)
	opts := strings.Split(uri[:p], ";")
	data := []byte(uri[p+1:])

	// Extract MIME type (first option)
	mime := opts[0]
	if mime == "" {
		mime = "application/octet-stream"
	}

	// Check if data is base64 encoded (last option is "base64")
	if opts[len(opts)-1] == "base64" {
		// Perform base64 decoding
		data = bytes.TrimRight(data, "=")
		result := make([]byte, base64.RawStdEncoding.DecodedLen(len(data)))
		n, err := base64.RawStdEncoding.Decode(result, data)
		if err != nil {
			return nil, "", fmt.Errorf("base64 decode failed: %w", err)
		}
		data = result[:n]
	} else {
		// URL-decode the data if not base64 encoded
		decoded, err := url.QueryUnescape(string(data))
		if err != nil {
			return nil, "", fmt.Errorf("URL decode failed: %w", err)
		}
		data = []byte(decoded)
	}

	return data, mime, nil
}

// ParseDataUri is an alias for ParseDataURI to maintain backward compatibility.
//
// Deprecated: Use ParseDataURI instead.
func ParseDataUri(uri string) ([]byte, string, error) {
	return ParseDataURI(uri)
}
