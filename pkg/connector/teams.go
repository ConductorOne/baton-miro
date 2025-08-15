package connector

import (
	"context"
	"fmt"
	"strings"

	"github.com/conductorone/baton-miro/pkg/miro"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	grant "github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

type teamBuilder struct {
	resourceType   *v2.ResourceType
	client         *miro.Client
	organizationId string
}

const (
	nonTeamTeamRole   = "non_team"
	memberTeamRole    = "member"
	adminTeamRole     = "admin"
	teamGuestTeamRole = "team_guest"
)

var teamRoles = []string{
	nonTeamTeamRole,
	memberTeamRole,
	adminTeamRole,
	teamGuestTeamRole,
}

// ResourceType returns the resource type for the team builder.
func (o *teamBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return teamResourceType
}

func teamResource(team *miro.Team) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"name": team.Name,
		"id":   team.Id,
	}

	teamTraitOptions := []rs.GroupTraitOption{
		rs.WithGroupProfile(profile),
	}
	resource, err := rs.NewGroupResource(team.Name, teamResourceType, team.Id, teamTraitOptions)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// newTeamBuilder creates a new team builder.
func newTeamBuilder(client *miro.Client, organizationId string) *teamBuilder {
	return &teamBuilder{
		resourceType:   teamResourceType,
		client:         client,
		organizationId: organizationId,
	}
}

// List returns the teams for an organization.
func (g *teamBuilder) List(ctx context.Context, _ *v2.ResourceId, pagination *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	bag, cursor, err := parsePageToken(pagination.Token, &v2.ResourceId{ResourceType: g.resourceType.Id})
	if err != nil {
		return nil, "", nil, wrapError(err, "failed to parse page token")
	}

	response, annos, err := g.client.GetTeams(ctx, g.organizationId, cursor, resourcePageSize)
	if err != nil {
		return nil, "", annos, wrapError(err, "failed to get teams")
	}

	var resources []*v2.Resource
	for _, team := range response.Data {
		team := team
		resource, err := teamResource(&team)
		if err != nil {
			return nil, "", annos, wrapError(err, "failed to create team resource")
		}

		resources = append(resources, resource)
	}

	if response.Cursor == "" {
		return resources, "", annos, nil
	}

	nextCursor, err := handleNextPage(bag, response.Cursor)
	if err != nil {
		return nil, "", nil, wrapError(err, "failed to create next page cursor")
	}

	return resources, nextCursor, nil, nil
}

// Entitlements returns the entitlements for a team.
func (o *teamBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var rv []*v2.Entitlement

	for _, role := range teamRoles {
		assigmentOptions := []ent.EntitlementOption{
			ent.WithGrantableTo(userResourceType),
			ent.WithDescription(fmt.Sprintf("Has %s team role", resource.DisplayName)),
			ent.WithDisplayName(fmt.Sprintf("%s team role %s", resource.DisplayName, role)),
		}

		entitlement := ent.NewAssignmentEntitlement(resource, role, assigmentOptions...)
		rv = append(rv, entitlement)
	}

	return rv, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (o *teamBuilder) Grants(ctx context.Context, resource *v2.Resource, pagination *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	bag, cursor, err := parsePageToken(pagination.Token, &v2.ResourceId{ResourceType: o.resourceType.Id})
	if err != nil {
		return nil, "", nil, wrapError(err, "failed to parse page token")
	}

	response, annos, err := o.client.GetTeamMembers(ctx, o.organizationId, resource.Id.Resource, cursor, resourcePageSize)
	if err != nil {
		return nil, "", annos, wrapError(err, "failed to get team members")
	}
	var grants []*v2.Grant
	for _, member := range response.Data {
		if !contains(teamRoles, member.Role) {
			return nil, "", annos, wrapError(nil, "user does not have a valid team role")
		}

		userResourceId := &v2.ResourceId{
			ResourceType: userResourceType.Id,
			Resource:     member.Id,
		}

		g := grant.NewGrant(resource, member.Role, userResourceId)
		grants = append(grants, g)
	}

	if response.Cursor == "" {
		return grants, "", annos, nil
	}

	nextCursor, err := handleNextPage(bag, response.Cursor)
	if err != nil {
		return nil, "", annos, wrapError(err, "failed to create next page cursor")
	}

	return grants, nextCursor, nil, nil
}

// Grant invites a user to a team.
func (o *teamBuilder) Grant(ctx context.Context, principial *v2.Resource, entitlement *v2.Entitlement) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)

	if principial.Id.ResourceType != userResourceType.Id {
		err := fmt.Errorf("baton-miro: only users can be invated to team")

		l.Warn(
			err.Error(),
			zap.String("principal_id", principial.Id.Resource),
			zap.String("principal_type", principial.Id.ResourceType),
		)

		return nil, err
	}

	role, err := parseTeamRoleFromEntitlementID(entitlement.Id)
	if err != nil {
		return nil, err
	}
	if !contains(teamRoles, role) {
		err := fmt.Errorf("baton-miro: invalid team role %s", role)

		l.Warn(
			err.Error(),
			zap.String("role", role),
		)

		return nil, err
	}

	user, annos, err := o.client.GetOrganizationMember(ctx, o.organizationId, principial.Id.Resource)
	if err != nil {
		return annos, wrapError(err, "failed to get user")
	}

	_, annos, err = o.client.InviteTeamMember(ctx, o.organizationId, entitlement.Resource.Id.Resource, user.Email, role)
	if err != nil {
		return annos, wrapError(err, "failed to invite user to team")
	}

	return annos, nil
}

// Revoke removes a user from a team.
func (g *teamBuilder) Revoke(ctx context.Context, grant *v2.Grant) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)

	entitlement := grant.Entitlement
	principal := grant.Principal

	if principal.Id.ResourceType != userResourceType.Id {
		err := fmt.Errorf("baton-miro: only users can be revoked from team")

		l.Warn(
			err.Error(),
			zap.String("principal_id", principal.Id.Resource),
			zap.String("principal_type", principal.Id.ResourceType),
		)

		return nil, err
	}

	_, err := g.client.RemoveTeamMember(ctx, g.organizationId, entitlement.Resource.Id.Resource, principal.Id.Resource)
	if err != nil {
		return nil, wrapError(err, "failed to remove user from team")
	}

	return nil, nil
}

func parseTeamRoleFromEntitlementID(entitlementID string) (string, error) {
	parts := strings.Split(entitlementID, ":")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid entitlement ID: %s", entitlementID)
	}

	role := parts[2]
	return role, nil
}
