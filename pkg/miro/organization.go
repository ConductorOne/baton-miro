package miro

import (
	"context"
	"net/http"
	"net/url"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

type (
	Organization struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	}
	GetOrganizationMembersResponse struct {
		Limit  int32  `json:"limit"`
		Size   int32  `json:"size"`
		Cursor string `json:"cursor"`
		Data   []User `json:"data"`
	}
)

func (c *Client) GetOrganizationMembers(ctx context.Context, organizationId, cursor string, limit int32, query ...queryFunction) (*GetOrganizationMembersResponse, *http.Response, error) {
	stringUrl, err := url.JoinPath(c.baseUrl, "v2/orgs", organizationId, "members")
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

	query = append(query, WithLimit(limit))
	if cursor != "" {
		query = append(query, WithCursor(cursor))
	}
	addQueryParams(req, query...)

	users := new(GetOrganizationMembersResponse)
	resp, err := c.Do(req, uhttp.WithJSONResponse(users))
	if err != nil {
		return nil, resp, err
	}

	return users, resp, nil
}

func (c *Client) GetOrganizationMember(ctx context.Context, organizationId, userId string) (*User, *http.Response, error) {
	stringValue, err := url.JoinPath(c.baseUrl, "v2/orgs", organizationId, "members", userId)
	if err != nil {
		return nil, nil, err
	}

	u, err := url.Parse(stringValue)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.NewRequest(ctx, http.MethodGet, u)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := c.Do(req, uhttp.WithJSONResponse(user))
	if err != nil {
		return nil, resp, err
	}

	return user, resp, nil
}
