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

const (
	OrganizationMembersUrl = "v2/orgs/%s/members"
)

func (c *Client) GetOrganizationMembers(ctx context.Context, organizationId string, cursor string, limit int32, query ...queryFunction) (*GetOrganizationMembersResponse, *http.Response, error) {
	u, err := buildResourceURL(c.baseUrl, fmt.Sprintf(OrganizationMembersUrl, organizationId))
	if err != nil {
		return nil, nil, err
	}

	query = append(query, WithLimit(limit))
	if cursor != "" {
		query = append(query, WithCursor(cursor))
	}
	addQueryParams(u, query...)

	var users GetOrganizationMembersResponse
	resp, err := c.doRequest(ctx, u, http.MethodGet, &users, nil)
	if err != nil {
		return nil, resp, err
	}

	return &users, resp, nil
}

func (c *Client) GetOrganizationMember(ctx context.Context, organizationId string, userId string) (*User, *http.Response, error) {
	u, err := buildResourceURL(c.baseUrl, fmt.Sprintf(OrganizationMembersUrl, organizationId), userId)
	if err != nil {
		return nil, nil, err
	}

	var user User
	resp, err := c.doRequest(ctx, u, http.MethodGet, &user, nil)
	if err != nil {
		return nil, resp, err
	}

	return &user, resp, nil
}
