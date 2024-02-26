package miro

import (
	"context"
	"net/http"
	"net/url"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

type Context struct {
	Type         string        `json:"type"`
	Team         *Team         `json:"team"`
	Scopes       []string      `json:"scopes"`
	User         *User         `json:"user"`
	Organization *Organization `json:"organization"`
}

func (c *Client) GetContext(ctx context.Context) (*Context, *http.Response, error) {
	stringUrl, err := url.JoinPath(c.baseUrl, "/v1/oauth-token")
	if err != nil {
		return nil, nil, err
	}

	u, err := url.Parse(stringUrl)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.NewRequest(ctx, http.MethodGet, u)
	if err != nil {
		return nil, nil, err
	}

	context := new(Context)
	resp, err := c.Do(req, uhttp.WithJSONResponse(context))
	if err != nil {
		return nil, resp, err
	}

	return context, resp, nil
}
