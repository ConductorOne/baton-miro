package miro

import (
	"context"
	"net/http"
	"net/url"
)

type Context struct {
	Type         string        `json:"type"`
	Team         *Team         `json:"team"`
	Scopes       []string      `json:"scopes"`
	User         *User         `json:"user"`
	Organization *Organization `json:"organization"`
}

const (
	accessTokenUrl = "/v1/oauth-token" //nolint:gosec // This is a URL path, not a hardcoded credential.
)

func (c *Client) GetContext(ctx context.Context) (*Context, *http.Response, error) {
	stringUrl, err := url.JoinPath(c.baseUrl, accessTokenUrl)
	if err != nil {
		return nil, nil, err
	}

	u, err := url.Parse(stringUrl)
	if err != nil {
		return nil, nil, err
	}

	accessToken := new(Context)
	resp, err := c.doRequest(ctx, u, http.MethodGet, accessToken, nil)
	if err != nil {
		return nil, resp, err
	}

	return accessToken, resp, nil
}
