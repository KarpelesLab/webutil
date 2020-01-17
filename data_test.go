package webutil_test

import (
	"bytes"
	"testing"

	"github.com/KarpelesLab/webutil"
)

func TestParseData(t *testing.T) {
	var cases = []struct {
		A    string
		R    []byte
		Mime string
	}{
		{"data:,", []byte{}, "application/octet-stream"},
		{"data:text/plain,Hello world", []byte("Hello world"), "text/plain"},
		{"data:text/plain,Hello+world", []byte("Hello world"), "text/plain"},
		{"data:text/plain,Hello%20world", []byte("Hello world"), "text/plain"},
		{"data:;base64,Zm9v", []byte("foo"), "application/octet-stream"},
		{"data:;base64,Zm9vYg", []byte("foob"), "application/octet-stream"},
		{"data://text/plain;base64,SSBsb3ZlIFBIUAo=", []byte("I love PHP\n"), "text/plain"}, // from php doc
	}

	for _, c := range cases {
		b, mime, err := webutil.ParseDataUri(c.A)

		if err != nil {
			t.Errorf("test failed, error %s", err)
		}

		if bytes.Compare(c.R, b) != 0 {
			t.Errorf("test failed, %s result %s, expected %v", c.A, b, c.R)
		}

		if c.Mime != mime {
			t.Errorf("test failed, %s mime result %s, expected %s", c.A, mime, c.Mime)
		}
	}
}
