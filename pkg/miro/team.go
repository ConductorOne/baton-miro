package miro

import (
	"context"
	"fmt"
	"net/http"
)

type (
	Team struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	}
	GetTeamsResponse struct {
		Limit  int32  `json:"limit"`
		Size   int32  `json:"size"`
		Cursor string `json:"cursor"`
		Data   []Team `json:"data"`
	}
	TeamMember struct {
		Id   string `json:"id"`
		Role string `json:"role"`
	}
	GetTeamMembersResponse struct {
		Limit  int32        `json:"limit"`
		Size   int32        `json:"size"`
		Cursor string       `json:"cursor"`
		Data   []TeamMember `json:"data"`
	}
	InviteTeamMemberBody struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	InviteTeamMemberResponse struct {
		TeamId string `json:"teamId"`
		Role   string `json:"role"`
		UserId string `json:"id"`
	}
)

func (c *Client) GetTeams(ctx context.Context, organizationId, cursor string, limit int32, query ...queryFunction) (*GetTeamsResponse, *http.Response, error) {
	url := fmt.Sprintf("%s/v2/orgs/%s/teams", c.baseUrl, organizationId)

	req, err := c.newRequestWithDefaultHeaders(ctx, http.MethodGet, url)
	if err != nil {
		return nil, nil, err
	}

	query = append(query, WithCursor(cursor), WithLimit(limit))
	addQueryParams(req, query...)

	teams := new(GetTeamsResponse)
	resp, err := c.do(req, teams)
	if err != nil {
		return nil, resp, err
	}

	return teams, resp, nil
}

func (c *Client) GetTeamMembers(ctx context.Context, organizationId, teamId, cursor string, limit int32, query ...queryFunction) (*GetTeamMembersResponse, *http.Response, error) {
	url := fmt.Sprintf("%s/v2/orgs/%s/teams/%s/member", c.baseUrl, organizationId, teamId)

	req, err := c.newRequestWithDefaultHeaders(ctx, http.MethodGet, url)
	if err != nil {
		return nil, nil, err
	}

	query = append(query, WithCursor(cursor), WithLimit(limit))
	addQueryParams(req, query...)

	teamMembers := new(GetTeamMembersResponse)
	resp, err := c.do(req, teamMembers)
	if err != nil {
		return nil, resp, err
	}

	return teamMembers, resp, nil
}

func (c *Client) InviteTeamMember(ctx context.Context, organizationId, teamId, email, role string) (*InviteTeamMemberResponse, *http.Response, error) {
	url := fmt.Sprintf("%s/v2/orgs/%s/teams/%s/member", c.baseUrl, organizationId, teamId)

	req, err := c.newRequestWithDefaultHeaders(ctx, http.MethodPost, url, InviteTeamMemberBody{
		Email: email,
		Role:  role,
	})
	if err != nil {
		return nil, nil, err
	}

	inviteTeamMemberResponse := new(InviteTeamMemberResponse)
	resp, err := c.do(req, inviteTeamMemberResponse)
	if err != nil {
		return nil, resp, err
	}

	return inviteTeamMemberResponse, resp, nil
}
