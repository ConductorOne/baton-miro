package miro

import (
	"context"
	"fmt"
	"net/http"
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
	url := fmt.Sprintf("%s/v1/orgs/%s/members", c.baseUrl, organizationId)

	req, err := c.newRequestWithDefaultHeaders(ctx, http.MethodGet, url)
	if err != nil {
		return nil, nil, err
	}

	query = append(query, WithCursor(cursor), WithLimit(limit))
	addQueryParams(req, query...)

	users := new(GetOrganizationMembersResponse)
	resp, err := c.do(req, users)
	if err != nil {
		return nil, resp, err
	}

	return users, resp, nil
}

func (c *Client) GetOrganizationMember(ctx context.Context, organizationId, userId string) (*User, *http.Response, error) {
	url := fmt.Sprintf("%s/v1/orgs/%s/members/%s", c.baseUrl, organizationId, userId)

	req, err := c.newRequestWithDefaultHeaders(ctx, http.MethodGet, url)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := c.do(req, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, nil
}
