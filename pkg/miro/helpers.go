package miro

import (
	"fmt"
	"net/url"
)

func buildResourceURL(baseURL string, endpoint string, elems ...string) (*url.URL, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	joined, err := url.JoinPath(u.Path, append([]string{endpoint}, elems...)...)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	u.Path = joined

	return u, nil
}
