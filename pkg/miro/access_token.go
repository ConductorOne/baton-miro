package miro

import (
	"context"
	"net/http"

	"github.com/conductorone/baton-sdk/pkg/annotations"
)

// Context is the context for the Miro client.
type Context struct {
	Type         string        `json:"type"`
	Team         *Team         `json:"team"`
	Scopes       []string      `json:"scopes"`
	User         *User         `json:"user"`
	Organization *Organization `json:"organization"`
}

// accessTokenUrl is the URL for the Access Token endpoint.
const (
	accessTokenUrl = "/v1/oauth-token" //nolint:gosec // This is a URL path, not a hardcoded credential.
)

// GetContext gets the context for the Miro client.
func (c *Client) GetContext(ctx context.Context) (*Context, annotations.Annotations, error) {
	accessToken := new(Context)
	_, annos, err := c.doRequest(ctx, accessTokenUrl, http.MethodGet, accessToken, nil)
	if err != nil {
		return nil, annos, err
	}

	return accessToken, annos, nil
}
