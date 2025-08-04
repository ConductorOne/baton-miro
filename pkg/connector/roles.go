package connector

import (
	"context"
	"fmt"
	"strings"

	"github.com/conductorone/baton-miro/pkg/miro"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
)

const (
	defaultRoleKey = "ORGANIZATION_INTERNAL_USER"
	assignedRole   = "assigned"
)

var organizationRoles = []string{
	assignedRole,
}

// roleDefinition is the definition of a role.
type roleDefinition struct {
	ID          string
	DisplayName string
	RoleKey     string
}

// roleDefinitions is the map of role definitions keyed by ID.
var roleDefinitions = map[string]roleDefinition{
	"organization_internal_admin":  {ID: "organization_internal_admin", DisplayName: "Organization Admin", RoleKey: "ORGANIZATION_INTERNAL_ADMIN"},
	"organization_internal_user":   {ID: "organization_internal_user", DisplayName: "Organization Internal User", RoleKey: "ORGANIZATION_INTERNAL_USER"},
	"organization_external_user":   {ID: "organization_external_user", DisplayName: "Organization External User", RoleKey: "ORGANIZATION_EXTERNAL_USER"},
	"organization_team_guest_user": {ID: "organization_team_guest_user", DisplayName: "Team Guest User", RoleKey: "ORGANIZATION_TEAM_GUEST_USER"},
}

// roleBuilder is the builder for the role resource type.
type roleBuilder struct {
	client         *miro.Client
	resourceType   *v2.ResourceType
	organizationId string
}

// ResourceType returns the resource type for the role builder.
func (r *roleBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return r.resourceType
}

// List returns the resources for the role builder.
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

// Entitlements returns the entitlements for the role builder.
func (r *roleBuilder) Entitlements(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var rv []*v2.Entitlement

	for _, role := range organizationRoles {
		assigmentOptions := []entitlement.EntitlementOption{
			entitlement.WithGrantableTo(userResourceType),
			entitlement.WithDescription(fmt.Sprintf("Has %s organization role", resource.DisplayName)),
			entitlement.WithDisplayName(fmt.Sprintf("%s organization role %s", resource.DisplayName, role)),
		}

		entitlement := entitlement.NewAssignmentEntitlement(resource, role, assigmentOptions...)
		rv = append(rv, entitlement)
	}

	return rv, "", nil, nil
}

// Grants returns empty grants - role grants are now emitted from user resources.
func (r *roleBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return []*v2.Grant{}, "", nil, nil
}

// Grant grants a role to a principal.
func (r *roleBuilder) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) ([]*v2.Grant, annotations.Annotations, error) {
	userID := principal.Id.Resource
	roleID := entitlement.Resource.Id.Resource

	roleDefinition, ok := roleDefinitions[roleID]
	if !ok {
		return nil, nil, fmt.Errorf("role key not found for ID: %s", roleID)
	}
	roleKey := roleDefinition.RoleKey

	scimUser, annos, err := r.client.GetUser(ctx, userID)
	if err != nil {
		return nil, annos, fmt.Errorf("failed to get user %s: %w", userID, err)
	}

	for _, userRole := range scimUser.Roles {
		if userRole.Value == roleKey {
			return nil, annotations.New(&v2.GrantAlreadyExists{}), nil
		}
	}

	_, annos, err = r.client.UpdateUserRole(ctx, userID, roleKey)
	if err != nil {
		return nil, annos, fmt.Errorf("failed to update user role for user %s: %w", userID, err)
	}

	g := grant.NewGrant(entitlement.Resource, assignedRole, principal.Id)
	return []*v2.Grant{g}, annos, nil
}

// Revoke revokes a role from a principal.
func (r *roleBuilder) Revoke(ctx context.Context, g *v2.Grant) (annotations.Annotations, error) {
	userID := g.Principal.Id.Resource
	roleID := g.Entitlement.Resource.Id.Resource

	roleDefinition, ok := roleDefinitions[roleID]
	if !ok {
		return nil, fmt.Errorf("role key not found for slug: %s", roleID)
	}
	roleToRevokeKey := roleDefinition.RoleKey

	user, annos, err := r.client.GetUser(ctx, userID)
	if err != nil {
		return annos, fmt.Errorf("failed to get user %s: %w", userID, err)
	}

	isRoleAlreadyRevoked := true
	for _, userRole := range user.Roles {
		if strings.EqualFold(userRole.Value, roleToRevokeKey) {
			isRoleAlreadyRevoked = false
			break
		}
	}

	if isRoleAlreadyRevoked {
		return annotations.New(&v2.GrantAlreadyRevoked{}), nil
	}

	if strings.EqualFold(roleToRevokeKey, defaultRoleKey) {
		return annotations.New(&v2.GrantAlreadyRevoked{}), nil
	}

	_, annos, err = r.client.UpdateUserRole(ctx, userID, defaultRoleKey)
	if err != nil {
		return annos, fmt.Errorf("failed to set default role for user %s: %w", userID, err)
	}

	return annos, nil
}

// newRoleBuilder creates a new role builder.
func newRoleBuilder(client *miro.Client, organizationId string) *roleBuilder {
	return &roleBuilder{
		client:         client,
		resourceType:   roleResourceType,
		organizationId: organizationId,
	}
}
