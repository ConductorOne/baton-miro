package miro

import (
	"context"
	"fmt"
	"net/http"
)

type Context struct {
	Type         string        `json:"type"`
	Team         *Team         `json:"team"`
	Scopes       []string      `json:"scopes"`
	User         *User         `json:"user"`
	Organization *Organization `json:"organization"`
}

func (c *Client) GetContext(ctx context.Context) (*Context, *http.Response, error) {
	url := fmt.Sprint(c.baseUrl, "/v1/oauth-token")

	req, err := c.newRequestWithDefaultHeaders(ctx, http.MethodGet, url)
	if err != nil {
		return nil, nil, err
	}

	context := new(Context)
	resp, err := c.do(req, context)
	if err != nil {
		return nil, resp, err
	}

	return context, resp, nil
}
