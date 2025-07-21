package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-miro/pkg/miro"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
)

type roleDefinition struct {
	ID          string
	DisplayName string
	RoleKey     string
}

var roleDefinitions = []roleDefinition{
	{ID: "organization_internal_admin", DisplayName: "Organization Admin", RoleKey: "ORGANIZATION_INTERNAL_ADMIN"},
	{ID: "organization_internal_user", DisplayName: "Organization Internal User", RoleKey: "ORGANIZATION_INTERNAL_USER"},
	{ID: "organization_external_user", DisplayName: "Organization External User", RoleKey: "ORGANIZATION_EXTERNAL_USER"},
	{ID: "organization_team_guest_user", DisplayName: "Team Guest User", RoleKey: "ORGANIZATION_TEAM_GUEST_USER"},
}

type roleBuilder struct {
	client       *miro.Client
	resourceType *v2.ResourceType
}

func (r *roleBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return r.resourceType
}

func (r *roleBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var resources []*v2.Resource
	for _, role := range roleDefinitions {
		profile := map[string]interface{}{
			"role_id":   role.ID,
			"role_name": role.DisplayName,
		}
		roleResource, err := resource.NewRoleResource(
			role.DisplayName,
			r.resourceType,
			role.ID,
			[]resource.RoleTraitOption{resource.WithRoleProfile(profile)},
		)
		if err != nil {
			return nil, "", nil, fmt.Errorf("failed to create role resource: %w", err)
		}
		resources = append(resources, roleResource)
	}
	return resources, "", nil, nil
}

func (r *roleBuilder) Entitlements(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	assigmentOptions := []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDescription(fmt.Sprintf("Assigned to %s role", resource.DisplayName)),
		entitlement.WithDisplayName(fmt.Sprintf("%s role Assignment", resource.DisplayName)),
	}

	ent := entitlement.NewAssignmentEntitlement(
		resource,
		"assigned",
		assigmentOptions...,
	)

	return []*v2.Entitlement{ent}, "", nil, nil
}

func (r *roleBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (r *roleBuilder) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) ([]*v2.Grant, annotations.Annotations, error) {
	userID := principal.Id.Resource
	roleID := entitlement.Resource.Id.Resource

	var roleKey string
	for _, definition := range roleDefinitions {
		if definition.ID == roleID {
			roleKey = definition.RoleKey
			break
		}
	}

	if roleKey == "" {
		return nil, nil, fmt.Errorf("role key not found for ID: %s", roleID)
	}

	user, _, err := r.client.GetUser(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user %s: %w", userID, err)
	}

	for _, userRole := range user.Roles {
		if userRole.Value == roleKey {
			return nil, annotations.New(&v2.GrantAlreadyExists{}), nil
		}
	}

	_, _, err = r.client.UpdateUserRole(ctx, userID, roleKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update user role for user %s: %w", userID, err)
	}

	g := grant.NewGrant(entitlement.Resource, "assigned", principal.Id)
	return []*v2.Grant{g}, nil, nil
}

func (r *roleBuilder) Revoke(ctx context.Context, g *v2.Grant) (annotations.Annotations, error) {
	userID := g.Principal.Id.Resource
	entitlementID := g.Entitlement.Resource.Id.Resource

	var roleToRevokeKey string
	for _, definition := range roleDefinitions {
		if definition.ID == entitlementID {
			roleToRevokeKey = definition.RoleKey
			break
		}
	}
	if roleToRevokeKey == "" {
		return nil, fmt.Errorf("role key not found for ID: %s", entitlementID)
	}

	user, _, err := r.client.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user %s: %w", userID, err)
	}

	defaultRoleKey := "ORGANIZATION_INTERNAL_USER"
	currentUserRole := ""
	if len(user.Roles) > 0 {
		currentUserRole = user.Roles[0].Value
	}

	if currentUserRole != roleToRevokeKey {
		return annotations.New(&v2.GrantAlreadyRevoked{}), nil
	}

	if currentUserRole == defaultRoleKey {
		return annotations.New(&v2.GrantAlreadyRevoked{}), nil
	}

	_, _, err = r.client.UpdateUserRole(ctx, userID, defaultRoleKey)
	if err != nil {
		return nil, fmt.Errorf("failed to set default role for user %s: %w", userID, err)
	}

	return nil, nil
}

func newRoleBuilder(client *miro.Client) *roleBuilder {
	return &roleBuilder{
		client:       client,
		resourceType: roleResourceType,
	}
}
