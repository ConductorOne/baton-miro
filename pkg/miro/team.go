package miro

import (
	"context"
	"fmt"
	"net/http"

	"github.com/conductorone/baton-sdk/pkg/annotations"
)

type (
	// Team is the response from the GetTeams endpoint.
	Team struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	}
	// GetTeamsResponse is the response from the GetTeams endpoint.
	GetTeamsResponse struct {
		Limit  int32  `json:"limit"`
		Size   int32  `json:"size"`
		Cursor string `json:"cursor"`
		Data   []Team `json:"data"`
	}
	// TeamMember is the response from the GetTeamMembers endpoint.
	TeamMember struct {
		Id         string `json:"id"`
		Role       string `json:"role"`
		CreatedAt  string `json:"createdAt"`
		CreatedBy  string `json:"createdBy"`
		ModifiedAt string `json:"modifiedAt"`
		ModifiedBy string `json:"modifiedBy"`
		TeamId     string `json:"teamId"`
		Type       string `json:"type"`
	}
	// GetTeamMembersResponse is the response from the GetTeamMembers endpoint.
	GetTeamMembersResponse struct {
		Limit  int32        `json:"limit"`
		Size   int32        `json:"size"`
		Cursor string       `json:"cursor"`
		Data   []TeamMember `json:"data"`
		Type   string       `json:"type"`
	}
	// InviteTeamMemberBody is the body for the InviteTeamMember endpoint.
	InviteTeamMemberBody struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	// InviteTeamMemberResponse is the response from the InviteTeamMember endpoint.
	InviteTeamMemberResponse struct {
		TeamId string `json:"teamId"`
		Role   string `json:"role"`
		UserId string `json:"id"`
	}
)

const (
	TeamsUrl       = "/v2/orgs/%s/teams"
	TeamMembersUrl = "/v2/orgs/%s/teams/%s/members"
)

// GetTeams gets the teams for a given organization.
func (c *Client) GetTeams(ctx context.Context, organizationId string, cursor string, limit int32, opts ...ReqOpt) (*GetTeamsResponse, annotations.Annotations, error) {
	teamsUrl, err := buildResourceURL(fmt.Sprintf(TeamsUrl, organizationId))
	if err != nil {
		return nil, nil, err
	}

	requestOpts := []ReqOpt{WithLimit(limit)}
	if cursor != "" {
		requestOpts = append(requestOpts, WithCursor(cursor))
	}
	requestOpts = append(requestOpts, opts...)

	var teams GetTeamsResponse
	_, annos, err := c.doRequest(ctx, teamsUrl.String(), http.MethodGet, &teams, nil, requestOpts...)
	if err != nil {
		return nil, annos, err
	}

	return &teams, annos, nil
}

// GetTeamMembers gets the team members for a given organization and team.
func (c *Client) GetTeamMembers(ctx context.Context, organizationId string, teamId string, cursor string, limit int32, opts ...ReqOpt) (*GetTeamMembersResponse, annotations.Annotations, error) {
	teamMembersUrl, err := buildResourceURL(fmt.Sprintf(TeamMembersUrl, organizationId, teamId))
	if err != nil {
		return nil, nil, err
	}

	requestOpts := []ReqOpt{WithLimit(limit)}
	if cursor != "" {
		requestOpts = append(requestOpts, WithCursor(cursor))
	}
	requestOpts = append(requestOpts, opts...)

	var teamMembers GetTeamMembersResponse
	_, annos, err := c.doRequest(ctx, teamMembersUrl.String(), http.MethodGet, &teamMembers, nil, requestOpts...)
	if err != nil {
		return nil, annos, err
	}

	return &teamMembers, annos, nil
}

// InviteTeamMember invites a team member to a given organization and team.
func (c *Client) InviteTeamMember(ctx context.Context, organizationId string, teamId string, email string, role string) (*InviteTeamMemberResponse, annotations.Annotations, error) {
	teamMembersUrl, err := buildResourceURL(fmt.Sprintf(TeamMembersUrl, organizationId, teamId))
	if err != nil {
		return nil, nil, err
	}

	body := InviteTeamMemberBody{
		Email: email,
		Role:  role,
	}

	var inviteTeamMemberResponse InviteTeamMemberResponse
	_, annos, err := c.doRequest(ctx, teamMembersUrl.String(), http.MethodPost, &inviteTeamMemberResponse, body)
	if err != nil {
		return nil, annos, err
	}

	return &inviteTeamMemberResponse, annos, nil
}

// RemoveTeamMember removes a team member from a given organization and team.
func (c *Client) RemoveTeamMember(ctx context.Context, organizationId string, teamId string, userId string) (annotations.Annotations, error) {
	teamMembersUrl, err := buildResourceURL(fmt.Sprintf(TeamMembersUrl, organizationId, teamId), userId)
	if err != nil {
		return nil, err
	}

	_, annos, err := c.doRequest(ctx, teamMembersUrl.String(), http.MethodDelete, nil, nil)
	if err != nil {
		return annos, err
	}

	return annos, nil
}
