package miro

import (
	"context"
	"fmt"
	"net/http"

	"github.com/conductorone/baton-sdk/pkg/annotations"
)

type (
	// Organization is the response from the GetOrganization endpoint.
	Organization struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	}
	// GetOrganizationMembersResponse is the response from the GetOrganizationMembers endpoint.
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

// GetOrganizationMembers gets the organization members for a given organization.
func (c *Client) GetOrganizationMembers(ctx context.Context, organizationId string, cursor string, limit int32, opts ...ReqOpt) (*GetOrganizationMembersResponse, annotations.Annotations, error) {
	getOrganizationMembersUrl, err := buildResourceURL(fmt.Sprintf(OrganizationMembersUrl, organizationId))
	if err != nil {
		return nil, nil, err
	}

	requestOpts := []ReqOpt{WithLimit(limit)}
	if cursor != "" {
		requestOpts = append(requestOpts, WithCursor(cursor))
	}
	requestOpts = append(requestOpts, opts...)

	var users GetOrganizationMembersResponse
	_, annos, err := c.doRequest(ctx, getOrganizationMembersUrl.String(), http.MethodGet, &users, nil, requestOpts...)
	if err != nil {
		return nil, annos, err
	}

	return &users, annos, nil
}

// GetOrganizationMember gets the organization member for a given organization and user.
func (c *Client) GetOrganizationMember(ctx context.Context, organizationId string, userId string) (*User, annotations.Annotations, error) {
	getOrganizationMemberUrl, err := buildResourceURL(fmt.Sprintf(OrganizationMembersUrl, organizationId), userId)
	if err != nil {
		return nil, nil, err
	}

	var user User
	_, annos, err := c.doRequest(ctx, getOrganizationMemberUrl.String(), http.MethodGet, &user, nil)
	if err != nil {
		return nil, annos, err
	}

	return &user, annos, nil
}
