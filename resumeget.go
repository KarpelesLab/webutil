package webutil

import (
	"fmt"
	"io"
	"net/http"
)

type resumeGET struct {
	req    *http.Request
	resp   *http.Response
	client *http.Client
	pos    int64 // current position
	size   int64
}

// discardAndCloseBody closes a body we don't need, but if it seems small it'll discard it so the connection can be re-used
func discardAndCloseBody(resp *http.Response) {
	const maxBodySlurpSize = 2 << 10
	if resp.ContentLength == -1 || resp.ContentLength <= maxBodySlurpSize {
		io.CopyN(io.Discard, resp.Body, maxBodySlurpSize)
	}
	resp.Body.Close()
}

// Get will perform a GET on the specified URL and return an object
// that will read data and automatically resume the download if an error
// occurs. Additionally it will also check the HTTP Status and return
// an error if it is not successful.
//
// If the response does not contain a Content-Length or if the server does
// not support Range headers, this will just perform a regular GET and won't
// be able to resume anything.
func Get(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// DefaultClient will handle redirects for us
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case http.StatusOK, http.StatusNoContent:
		// all good
	default:
		discardAndCloseBody(resp)
		return nil, fmt.Errorf("HTTP response %s", resp.Status)
	}

	obj := &resumeGET{
		req:    resp.Request, // use resp.Request so we can skip all the redirects (if any) on resume
		resp:   resp,
		size:   resp.ContentLength,
		client: http.DefaultClient,
	}

	return obj, nil
}

func (r *resumeGET) Read(b []byte) (int, error) {
	if r.resp != nil {
		n, err := r.resp.Body.Read(b)
		if err == nil || (err != nil && n > 0) {
			// no error, or an error happened but we did manage to read some stuff, do not report the error yet
			r.pos += int64(n)
			return n, nil
		}

		// error + could not read anything
		if r.pos >= r.size {
			// end of file
			return 0, io.EOF
		}

		// error, close & discard r.resp
		r.resp.Body.Close()
		r.resp = nil
	}

	// we're not at EOF yet, let's try to resume this
	r.req.Header.Set("Range", fmt.Sprintf("bytes=%d-", r.pos))
	resp, err := r.client.Do(r.req)
	if err != nil {
		return 0, fmt.Errorf("while attempting to resume: %w", err)
	}
	if resp.StatusCode != http.StatusPartialContent {
		discardAndCloseBody(resp)
		return 0, fmt.Errorf("expected partial content, got %w", resp.Status)
	}

	// store resp
	r.resp = resp
	n, err := r.resp.Body.Read(b)
	r.pos += int64(n)
	return n, err
}

func (r *resumeGET) Close() error {
	return r.resp.Body.Close()
}
