package webutil_test

import (
	"bytes"
	"testing"

	"github.com/KarpelesLab/webutil"
)

func TestParseDataURI(t *testing.T) {
	testCases := []struct {
		input    string
		expected []byte
		mime     string
		name     string
	}{
		{
			name:     "Empty data",
			input:    "data:,",
			expected: []byte{},
			mime:     "application/octet-stream",
		},
		{
			name:     "Plain text without encoding",
			input:    "data:text/plain,Hello world",
			expected: []byte("Hello world"),
			mime:     "text/plain",
		},
		{
			name:     "Plain text with plus encoding",
			input:    "data:text/plain,Hello+world",
			expected: []byte("Hello world"),
			mime:     "text/plain",
		},
		{
			name:     "Plain text with percent encoding",
			input:    "data:text/plain,Hello%20world",
			expected: []byte("Hello world"),
			mime:     "text/plain",
		},
		{
			name:     "Base64 with default mime type",
			input:    "data:;base64,Zm9v",
			expected: []byte("foo"),
			mime:     "application/octet-stream",
		},
		{
			name:     "Base64 with non-multiple-of-four length",
			input:    "data:;base64,Zm9vYg",
			expected: []byte("foob"),
			mime:     "application/octet-stream",
		},
		{
			name:     "Base64 with slashes in scheme part",
			input:    "data://text/plain;base64,SSBsb3ZlIFBIUAo=",
			expected: []byte("I love PHP\n"),
			mime:     "text/plain",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, mime, err := webutil.ParseDataURI(tc.input)

			if err != nil {
				t.Fatalf("Failed to parse valid data URI: %v", err)
			}

			if !bytes.Equal(tc.expected, data) {
				t.Errorf("Data mismatch: got %q, want %q", data, tc.expected)
			}

			if tc.mime != mime {
				t.Errorf("MIME type mismatch: got %q, want %q", mime, tc.mime)
			}
		})
	}

	// Test error cases
	errorCases := []struct {
		input string
		name  string
	}{
		{
			name:  "Invalid prefix",
			input: "invalid:data",
		},
		{
			name:  "Missing comma separator",
			input: "data:text/plain",
		},
	}

	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := webutil.ParseDataURI(tc.input)
			if err == nil {
				t.Errorf("Expected error for invalid data URI: %s", tc.input)
			}
		})
	}
}
