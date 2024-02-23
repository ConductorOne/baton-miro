package miro

import (
	"context"
	"net/http"
	"net/url"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
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
	stringUrl, err := url.JoinPath(c.baseUrl, "v2/orgs", organizationId, "teams")
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

	if cursor != "" {
		query = append(query, WithCursor(cursor))
	}
	addQueryParams(req, query...)

	teams := new(GetTeamsResponse)
	resp, err := c.Do(req, uhttp.WithJSONResponse(teams))
	if err != nil {
		return nil, resp, err
	}

	return teams, resp, nil
}

func (c *Client) GetTeamMembers(ctx context.Context, organizationId, teamId, cursor string, limit int32, query ...queryFunction) (*GetTeamMembersResponse, *http.Response, error) {
	stringUrl, err := url.JoinPath(c.baseUrl, "v2/orgs", organizationId, "teams", teamId, "members")
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

	teamMembers := new(GetTeamMembersResponse)
	resp, err := c.Do(req, uhttp.WithJSONResponse(teamMembers))
	if err != nil {
		return nil, resp, err
	}

	return teamMembers, resp, nil
}

func (c *Client) InviteTeamMember(ctx context.Context, organizationId, teamId, email, role string) (*InviteTeamMemberResponse, *http.Response, error) {
	stringUrl, err := url.JoinPath(c.baseUrl, "v2/orgs", organizationId, "teams", teamId, "members")
	if err != nil {
		return nil, nil, err
	}

	u, err := url.Parse(stringUrl)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.NewRequest(ctx, http.MethodPost, u, uhttp.WithJSONBody(InviteTeamMemberBody{
		Email: email,
		Role:  role,
	}))
	if err != nil {
		return nil, nil, err
	}

	inviteTeamMemberResponse := new(InviteTeamMemberResponse)
	resp, err := c.Do(req, uhttp.WithJSONResponse(inviteTeamMemberResponse))
	if err != nil {
		return nil, resp, err
	}

	return inviteTeamMemberResponse, resp, nil
}

func (c *Client) RemoveTeamMember(ctx context.Context, organizationId, teamId, userId string) (*http.Response, error) {
	stringUrl, err := url.JoinPath(c.baseUrl, "v2/orgs", organizationId, "teams", teamId, "members", userId)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(stringUrl)
	if err != nil {
		return nil, err
	}

	req, err := c.NewRequest(ctx, http.MethodDelete, u)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
