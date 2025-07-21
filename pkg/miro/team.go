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

const (
	TeamsUrl       = "v2/orgs/%s/teams"
	TeamMembersUrl = "v2/orgs/%s/teams/%s/members"
)

func (c *Client) GetTeams(ctx context.Context, organizationId string, cursor string, limit int32, query ...queryFunction) (*GetTeamsResponse, *http.Response, error) {
	teamsUrl, err := buildResourceURL(c.baseUrl, fmt.Sprintf(TeamsUrl, organizationId))
	if err != nil {
		return nil, nil, err
	}

	if cursor != "" {
		query = append(query, WithCursor(cursor))
	}
	addQueryParams(teamsUrl, query...)

	var teams GetTeamsResponse
	resp, err := c.doRequest(ctx, teamsUrl, http.MethodGet, &teams, nil)
	if err != nil {
		return nil, resp, err
	}

	return &teams, resp, nil
}

func (c *Client) GetTeamMembers(ctx context.Context, organizationId string, teamId string, cursor string, limit int32, query ...queryFunction) (*GetTeamMembersResponse, *http.Response, error) {
	teamMembersUrl, err := buildResourceURL(c.baseUrl, fmt.Sprintf(TeamMembersUrl, organizationId, teamId))
	if err != nil {
		return nil, nil, err
	}

	query = append(query, WithLimit(limit))
	if cursor != "" {
		query = append(query, WithCursor(cursor))
	}
	addQueryParams(teamMembersUrl, query...)

	var teamMembers GetTeamMembersResponse
	resp, err := c.doRequest(ctx, teamMembersUrl, http.MethodGet, &teamMembers, nil)
	if err != nil {
		return nil, resp, err
	}

	return &teamMembers, resp, nil
}

func (c *Client) InviteTeamMember(ctx context.Context, organizationId string, teamId string, email string, role string) (*InviteTeamMemberResponse, *http.Response, error) {
	teamMembersUrl, err := buildResourceURL(c.baseUrl, fmt.Sprintf(TeamMembersUrl, organizationId, teamId))
	if err != nil {
		return nil, nil, err
	}

	body := InviteTeamMemberBody{
		Email: email,
		Role:  role,
	}

	var inviteTeamMemberResponse InviteTeamMemberResponse
	resp, err := c.doRequest(ctx, teamMembersUrl, http.MethodPost, &inviteTeamMemberResponse, body)
	if err != nil {
		return nil, resp, err
	}

	return &inviteTeamMemberResponse, resp, nil
}

func (c *Client) RemoveTeamMember(ctx context.Context, organizationId string, teamId string, userId string) (*http.Response, error) {
	teamMembersUrl, err := buildResourceURL(c.baseUrl, fmt.Sprintf(TeamMembersUrl, organizationId, teamId), userId)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, teamMembersUrl, http.MethodDelete, nil, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
