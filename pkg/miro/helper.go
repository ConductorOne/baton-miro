package miro

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// ReqOpt represents a request option that can be applied to an HTTP request.
type ReqOpt func(*http.Request) *http.Request

// WithQueryParam adds a query parameter to the request.
func WithQueryParam(key, value string) ReqOpt {
	return func(req *http.Request) *http.Request {
		q := req.URL.Query()
		q.Add(key, value)
		req.URL.RawQuery = q.Encode()
		return req
	}
}

// WithLimit adds a limit query parameter to the request.
func WithLimit(limit int32) ReqOpt {
	return WithQueryParam("limit", strconv.Itoa(int(limit)))
}

// WithCursor adds a cursor query parameter to the request.
func WithCursor(cursor string) ReqOpt {
	if cursor == "" {
		return func(req *http.Request) *http.Request { return req }
	}
	return WithQueryParam("cursor", cursor)
}

// buildResourceURL builds a resource URL from an endpoint and path elements.
func buildResourceURL(endpoint string, elems ...string) (*url.URL, error) {
	pathElements := append([]string{endpoint}, elems...)
	joined, err := url.JoinPath("", pathElements...)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	u := &url.URL{
		Path: joined,
	}

	return u, nil
}
