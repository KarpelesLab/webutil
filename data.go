package webutil

import (
	"bytes"
	"encoding/base64"
	"errors"
	"net/url"
	"strings"
)

// ParseDataUri will parse a given data: uri and return its data and mime type.
func ParseDataUri(u string) ([]byte, string, error) {
	if !strings.HasPrefix(u, "data:") {
		return nil, "", errors.New("not a data uri")
	}
	u = u[5:]
	for len(u) > 0 && u[0] == '/' {
		u = u[1:]
	}

	p := strings.IndexByte(u, ',')

	if p == -1 {
		return nil, "", errors.New("could not locate base64 encoded value")
	}

	opts := strings.Split(u[:p], ";") // first opt will be mime type, last will be base64 if base64
	dat := []byte(u[p+1:])            // could be base64 encoded, we'll see this later

	mime := opts[0]
	if mime == "" {
		mime = "application/octet-stream"
	}

	if opts[len(opts)-1] == "base64" {
		// perform base64 decoding
		dat = bytes.TrimRight(dat, "=")
		res := make([]byte, base64.RawStdEncoding.DecodedLen(len(dat)))
		n, err := base64.RawStdEncoding.Decode(res, dat)

		if err != nil {
			return nil, "", err
		}

		dat = res[:n]
	} else {
		tmp, err := url.QueryUnescape(string(dat))
		if err != nil {
			return nil, "", err
		}
		dat = []byte(tmp)
	}

	return dat, mime, nil
}
