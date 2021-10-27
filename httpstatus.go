package webutil

import (
	"io/fs"
	"net/http"
)

func HTTPStatus(err error) int {
	switch err {
	case fs.ErrNotExist:
		return http.StatusNotFound
	case fs.ErrPermission:
		return http.StatusForbidden
	}

	switch e := err.(type) {
	case HttpError:
		return int(e)
	case interface{ HTTPStatus() int }:
		return e.HTTPStatus()
	case interface{ Unwrap() error }:
		return HTTPStatus(e.Unwrap())
	default:
		return 0
	}
}
