package miro

import (
	"context"
	"net/http"
)

const (
	UsersUrl = "/Users"
)

// CreateUser creates a new user in Miro using the SCIM API.
func (c *Client) CreateUser(ctx context.Context, user *CreateUserRequest) (*User, *http.Response, error) {
	u, err := buildResourceURL(ScimBaseUrl, UsersUrl)
	if err != nil {
		return nil, nil, err
	}

	var userResponse User
	resp, err := c.doRequest(ctx, u, http.MethodPost, &userResponse, user)
	if err != nil {
		return nil, resp, err
	}

	return &userResponse, resp, nil
}

// GetUser fetches a user by ID using the SCIM API.
func (c *Client) GetUser(ctx context.Context, userId string) (*ScimUser, *http.Response, error) {
	u, err := buildResourceURL(ScimBaseUrl, UsersUrl, userId)
	if err != nil {
		return nil, nil, err
	}

	var userResponse ScimUser
	resp, err := c.doRequest(ctx, u, http.MethodGet, &userResponse, nil)
	if err != nil {
		return nil, resp, err
	}

	return &userResponse, resp, nil
}

// UpdateUserRole updates the role of a user in Miro using the SCIM API.
func (c *Client) UpdateUserRole(ctx context.Context, userId string, role string) (*ScimUser, *http.Response, error) {
	u, err := buildResourceURL(ScimBaseUrl, UsersUrl, userId)
	if err != nil {
		return nil, nil, err
	}

	patchData := PatchOp{
		Schemas: []string{"urn:ietf:params:scim:api:messages:2.0:PatchOp"},
		Operations: []PatchOpItem{
			{
				Op:    "Replace",
				Path:  "roles.value",
				Value: role,
			},
		},
	}

	var userResponse ScimUser
	resp, err := c.doRequest(ctx, u, http.MethodPatch, &userResponse, &patchData)
	if err != nil {
		return nil, resp, err
	}

	return &userResponse, resp, nil
}
