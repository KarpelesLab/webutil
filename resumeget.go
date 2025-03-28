package webutil

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// resumeGET implements an io.ReadCloser that automatically resumes downloads
// when connections are interrupted.
type resumeGET struct {
	req    *http.Request
	resp   *http.Response
	client *http.Client
	pos    int64 // current position in bytes
	size   int64 // total size in bytes (if known)
}

// discardAndCloseBody closes a response body properly, ensuring connection reuse
// by first discarding the body content if it's small enough.
func discardAndCloseBody(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}

	const maxBodySlurpSize = 2 << 10 // 2KB
	if resp.ContentLength == -1 || resp.ContentLength <= maxBodySlurpSize {
		// For small or unknown size bodies, read a bit to enable connection reuse
		_, _ = io.CopyN(io.Discard, resp.Body, maxBodySlurpSize)
	}
	resp.Body.Close()
}

// Get performs an HTTP GET request that can automatically resume downloads.
//
// It returns an io.ReadCloser that will:
// 1. Automatically resume the download if the connection is interrupted
// 2. Handle proper error reporting based on HTTP status codes
//
// Limitations:
// - If the server doesn't support Range headers, it can't resume
// - If Content-Length isn't provided, size tracking won't be accurate
//
// The returned io.ReadCloser implements transparent resuming by using
// Range headers when an HTTP connection fails mid-download.
func Get(url string) (io.ReadCloser, error) {
	// Use context for better control and potential timeout in the future
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// DefaultClient handles redirects for us
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("performing request: %w", err)
	}

	// Check if the status code indicates success
	switch resp.StatusCode {
	case http.StatusOK, http.StatusNoContent:
		// Success, continue
	default:
		// Error status, clean up and return an error
		discardAndCloseBody(resp)
		return nil, HTTPError(resp.StatusCode)
	}

	// Create a resumeGET object that can handle interrupted downloads
	getter := &resumeGET{
		// Use resp.Request to retain any redirects that occurred
		req:    resp.Request,
		resp:   resp,
		size:   resp.ContentLength,
		client: http.DefaultClient,
	}

	return getter, nil
}

// Read implements io.Reader, handling automatic resumption of interrupted downloads.
func (r *resumeGET) Read(b []byte) (int, error) {
	// If we have an active response, try to read from it
	if r.resp != nil {
		n, err := r.resp.Body.Read(b)

		// If we read data, update position and return, even if there was an error
		// This ensures we return the data we have before handling any error
		if n > 0 {
			r.pos += int64(n)
			return n, err
		}

		// Error with no data read: decide if we should resume or return EOF
		if err != nil {
			// If we've already read the entire content or size is unknown and err is EOF,
			// we're done
			if (r.size > 0 && r.pos >= r.size) || (r.size < 0 && err == io.EOF) {
				return 0, io.EOF
			}

			// Otherwise, close the current response to prepare for resumption
			r.resp.Body.Close()
			r.resp = nil
		}
	}

	// No active response or previous response had an error, attempt to resume
	return r.resumeDownload(b)
}

// resumeDownload attempts to resume an interrupted download
// using Range headers.
func (r *resumeGET) resumeDownload(b []byte) (int, error) {
	// Set Range header to resume from current position
	r.req.Header.Set("Range", fmt.Sprintf("bytes=%d-", r.pos))

	// Perform the request with the Range header
	resp, err := r.client.Do(r.req)
	if err != nil {
		return 0, fmt.Errorf("resuming download: %w", err)
	}

	// Server must respond with 206 Partial Content for a successful range request
	if resp.StatusCode != http.StatusPartialContent {
		discardAndCloseBody(resp)
		return 0, fmt.Errorf("expected 206 Partial Content, got %w", HTTPError(resp.StatusCode))
	}

	// Store the new response
	r.resp = resp

	// Read data from the new response
	n, err := r.resp.Body.Read(b)
	r.pos += int64(n)
	return n, err
}

// Close implements io.Closer, ensuring the response body is properly closed.
func (r *resumeGET) Close() error {
	if r.resp == nil || r.resp.Body == nil {
		return nil
	}
	return r.resp.Body.Close()
}
