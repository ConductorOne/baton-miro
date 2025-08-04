package connector

import (
	"context"
	"testing"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/pagination"
)

// TestRoleDefinitions tests the role definitions.
func TestRoleDefinitions(t *testing.T) {
	expectedRoles := map[string]roleDefinition{
		"organization_internal_admin":  {ID: "organization_internal_admin", DisplayName: "Organization Admin", RoleKey: "ORGANIZATION_INTERNAL_ADMIN"},
		"organization_internal_user":   {ID: "organization_internal_user", DisplayName: "Organization Internal User", RoleKey: "ORGANIZATION_INTERNAL_USER"},
		"organization_external_user":   {ID: "organization_external_user", DisplayName: "Organization External User", RoleKey: "ORGANIZATION_EXTERNAL_USER"},
		"organization_team_guest_user": {ID: "organization_team_guest_user", DisplayName: "Team Guest User", RoleKey: "ORGANIZATION_TEAM_GUEST_USER"},
	}

	if len(roleDefinitions) != len(expectedRoles) {
		t.Errorf("roleDefinitions length = %v, want %v", len(roleDefinitions), len(expectedRoles))
	}

	for id, expected := range expectedRoles {
		actual, exists := roleDefinitions[id]
		if !exists {
			t.Errorf("roleDefinitions[%s] does not exist", id)
			continue
		}
		if actual.ID != expected.ID {
			t.Errorf("roleDefinitions[%s].ID = %v, want %v", id, actual.ID, expected.ID)
		}
		if actual.DisplayName != expected.DisplayName {
			t.Errorf("roleDefinitions[%s].DisplayName = %v, want %v", id, actual.DisplayName, expected.DisplayName)
		}
		if actual.RoleKey != expected.RoleKey {
			t.Errorf("roleDefinitions[%s].RoleKey = %v, want %v", id, actual.RoleKey, expected.RoleKey)
		}
	}
}

// TestRoleBuilder_ResourceType tests the resource type for the role builder.
func TestRoleBuilder_ResourceType(t *testing.T) {
	builder := &roleBuilder{
		resourceType: roleResourceType,
	}

	result := builder.ResourceType(context.Background())

	if result != roleResourceType {
		t.Errorf("ResourceType() = %v, want %v", result, roleResourceType)
	}
}

// TestRoleBuilder_List tests the list method for the role builder.
func TestRoleBuilder_List(t *testing.T) {
	builder := &roleBuilder{
		resourceType: roleResourceType,
	}

	resources, nextPage, _, err := builder.List(context.Background(), &v2.ResourceId{}, &pagination.Token{})

	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if nextPage != "" {
		t.Errorf("List() nextPage = %v, want empty string", nextPage)
	}

	expectedCount := len(roleDefinitions)
	if len(resources) != expectedCount {
		t.Errorf("List() count = %v, want %v", len(resources), expectedCount)
	}

	for i, resource := range resources {
		expected, exists := roleDefinitions[resource.Id.Resource]
		if !exists {
			t.Errorf("List()[%d].Id.Resource = %v, role definition not found", i, resource.Id.Resource)
			continue
		}
		if resource.Id.Resource != expected.ID {
			t.Errorf("List()[%d].Id.Resource = %v, want %v", i, resource.Id.Resource, expected.ID)
		}
		if resource.DisplayName != expected.DisplayName {
			t.Errorf("List()[%d].DisplayName = %v, want %v", i, resource.DisplayName, expected.DisplayName)
		}
	}
}

// TestRoleBuilder_Entitlements tests the entitlements for a role.
func TestRoleBuilder_Entitlements(t *testing.T) {
	builder := &roleBuilder{
		resourceType: roleResourceType,
	}

	resource := &v2.Resource{
		Id: &v2.ResourceId{
			ResourceType: roleResourceType.Id,
			Resource:     "organization_internal_admin",
		},
		DisplayName: "Organization Admin",
	}

	entitlements, nextPage, _, err := builder.Entitlements(context.Background(), resource, &pagination.Token{})

	if err != nil {
		t.Fatalf("Entitlements() error = %v", err)
	}

	if nextPage != "" {
		t.Errorf("Entitlements() nextPage = %v, want empty string", nextPage)
	}

	expectedCount := 1
	if len(entitlements) != expectedCount {
		t.Errorf("Entitlements() count = %v, want %v", len(entitlements), expectedCount)
	}

	if len(entitlements) > 0 {
		entitlement := entitlements[0]
		if entitlement.Resource.Id.Resource != "organization_internal_admin" {
			t.Errorf("Entitlements()[0].Resource.Id.Resource = %v, want organization_internal_admin", entitlement.Resource.Id.Resource)
		}
		if entitlement.Slug != assignedRole {
			t.Errorf("Entitlements()[0].Slug = %v, want %v", entitlement.Slug, assignedRole)
		}
	}
}

// TestRoleBuilder_Grants_EmptyResult tests that Grants returns empty list.
func TestRoleBuilder_Grants_EmptyResult(t *testing.T) {
	builder := &roleBuilder{
		resourceType:   roleResourceType,
		organizationId: "test-org-id",
	}

	resource := &v2.Resource{
		Id: &v2.ResourceId{
			ResourceType: roleResourceType.Id,
			Resource:     "organization_internal_admin",
		},
	}

	grants, nextPage, _, err := builder.Grants(context.Background(), resource, &pagination.Token{})

	if err != nil {
		t.Fatalf("Grants() unexpected error: %v", err)
	}

	if nextPage != "" {
		t.Errorf("Grants() nextPage = %v, want empty string", nextPage)
	}

	if grants == nil {
		t.Errorf("Grants() = nil, want empty slice")
	}

	if len(grants) != 0 {
		t.Errorf("Grants() length = %v, want 0 (role grants are now emitted from user resources)", len(grants))
	}
}
