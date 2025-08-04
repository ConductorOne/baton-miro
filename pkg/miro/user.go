package miro

import (
	"context"
	"net/http"

	"github.com/conductorone/baton-sdk/pkg/annotations"
)

// UsersUrl is the URL for the Users endpoint.
const (
	UsersUrl = "/Users"
)

// CreateUser creates a new user in Miro using the SCIM API.
func (c *Client) CreateUser(ctx context.Context, email string, firstName string, lastName string) (*User, annotations.Annotations, error) {
	createUserUrl, err := buildResourceURL(UsersUrl)
	if err != nil {
		return nil, nil, err
	}

	createUserReq := &CreateUserRequest{
		Schemas:  []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
		UserName: email,
		Name: RequestName{
			GivenName:  firstName,
			FamilyName: lastName,
		},
	}

	var userResponse User
	_, annos, err := c.doScimRequest(ctx, createUserUrl.String(), http.MethodPost, &userResponse, createUserReq)
	if err != nil {
		return nil, annos, err
	}

	return &userResponse, annos, nil
}

// GetUser fetches a user by ID using the SCIM API.
func (c *Client) GetUser(ctx context.Context, userId string) (*ScimUser, annotations.Annotations, error) {
	getUserUrl, err := buildResourceURL(UsersUrl, userId)
	if err != nil {
		return nil, nil, err
	}

	var userResponse ScimUser
	_, annos, err := c.doScimRequest(ctx, getUserUrl.String(), http.MethodGet, &userResponse, nil)
	if err != nil {
		return nil, annos, err
	}

	return &userResponse, annos, nil
}

// ReplaceUser completely replaces a user using the SCIM PUT API.
func (c *Client) ReplaceUser(ctx context.Context, userId string, user *ScimUser) (*ScimUser, annotations.Annotations, error) {
	replaceUserUrl, err := buildResourceURL(UsersUrl, userId)
	if err != nil {
		return nil, nil, err
	}

	var userResponse ScimUser
	_, annos, err := c.doScimRequest(ctx, replaceUserUrl.String(), http.MethodPut, &userResponse, user)
	if err != nil {
		return nil, annos, err
	}

	return &userResponse, annos, nil
}

// UpdateUserRole updates the role of a user in Miro using the SCIM API.
func (c *Client) UpdateUserRole(ctx context.Context, userId string, role string) (*ScimUser, annotations.Annotations, error) {
	updateUserRoleUrl, err := buildResourceURL(UsersUrl, userId)
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
	_, annos, err := c.doScimRequest(ctx, updateUserRoleUrl.String(), http.MethodPatch, &userResponse, &patchData)
	if err != nil {
		return nil, annos, err
	}

	return &userResponse, annos, nil
}
