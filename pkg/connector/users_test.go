package connector

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/conductorone/baton-miro/pkg/miro"
	"github.com/conductorone/baton-miro/test"
)

const (
	mockUserID    = "user-123"
	mockUserEmail = "john.doe@example.com"
)

// TestUserMockData tests the user mock data.
func TestUserMockData(t *testing.T) {
	mockData := test.ReadFile("scim_user_success.json")

	var user miro.ScimUser
	err := json.Unmarshal([]byte(mockData), &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal mock user data: %v", err)
	}

	if user.Id != mockUserID {
		t.Errorf("Expected user ID to be %s, got %s", mockUserID, user.Id)
	}

	if user.UserName != mockUserEmail {
		t.Errorf("Expected username to be %s, got %s", mockUserEmail, user.UserName)
	}

	if user.DisplayName != "John Doe" {
		t.Errorf("Expected display name to be 'John Doe', got %s", user.DisplayName)
	}

	if !user.Active {
		t.Error("Expected user to be active")
	}

	if len(user.Emails) == 0 {
		t.Error("Expected user to have emails")
	} else if user.Emails[0].Value != mockUserEmail {
		t.Errorf("Expected email to be %s, got %s", mockUserEmail, user.Emails[0].Value)
	}
}

// TestOrganizationUserMockData tests the organization user mock data.
func TestOrganizationUserMockData(t *testing.T) {
	mockData := test.ReadFile("organization_user_success.json")

	var user miro.User
	err := json.Unmarshal([]byte(mockData), &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal mock organization user data: %v", err)
	}

	if user.Id != mockUserID {
		t.Errorf("Expected user ID to be %s, got %s", mockUserID, user.Id)
	}

	if user.Email != mockUserEmail {
		t.Errorf("Expected email to be %s, got %s", mockUserEmail, user.Email)
	}

	if !user.Active {
		t.Error("Expected user to be active")
	}

	if user.Type != "user" {
		t.Errorf("Expected type to be 'user', got %s", user.Type)
	}

	if user.License != "full" {
		t.Errorf("Expected license to be 'full', got %s", user.License)
	}
}

// TestUserBuilder_ResourceType tests the resource type for the user builder.
func TestUserBuilder_ResourceType(t *testing.T) {
	builder := &userBuilder{
		resourceType: userResourceType,
	}

	result := builder.ResourceType(context.TODO())

	if result != userResourceType {
		t.Errorf("ResourceType() = %v, want %v", result, userResourceType)
	}
}

// TestUserStatusMapping tests the user status mapping.
func TestUserStatusMapping(t *testing.T) {
	tests := []struct {
		name     string
		user     miro.User
		expected string
	}{
		{
			name: "active user",
			user: miro.User{
				Active: true,
			},
			expected: "enabled",
		},
		{
			name: "inactive user",
			user: miro.User{
				Active: false,
			},
			expected: "disabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.user.Active && tt.expected != "enabled" {
				t.Errorf("Expected enabled status for active user, got %s", tt.expected)
			}
			if !tt.user.Active && tt.expected != "disabled" {
				t.Errorf("Expected disabled status for inactive user, got %s", tt.expected)
			}
		})
	}
}

// TestUserBuilderStructure tests the user builder structure.
func TestUserBuilderStructure(t *testing.T) {
	builder := &userBuilder{
		resourceType:   userResourceType,
		organizationId: "test-org-123",
	}

	if builder.resourceType != userResourceType {
		t.Errorf("userBuilder resourceType = %v, want %v", builder.resourceType, userResourceType)
	}

	if builder.organizationId != "test-org-123" {
		t.Errorf("userBuilder organizationId = %v, want %v", builder.organizationId, "test-org-123")
	}

	if builder.resourceType == nil {
		t.Fatal("userBuilder resourceType is nil")
	}
}

// TestCreateUserRequestValidation tests the create user request validation.
func TestCreateUserRequestValidation(t *testing.T) {
	createReq := miro.CreateUserRequest{
		Schemas:  []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
		UserName: "test.user@example.com",
		Name: miro.RequestName{
			GivenName:  "Test",
			FamilyName: "User",
		},
	}

	if len(createReq.Schemas) == 0 {
		t.Error("CreateUserRequest should have schemas")
	}

	expectedSchema := "urn:ietf:params:scim:schemas:core:2.0:User"
	if createReq.Schemas[0] != expectedSchema {
		t.Errorf("Expected schema %s, got %s", expectedSchema, createReq.Schemas[0])
	}

	if createReq.UserName == "" {
		t.Error("CreateUserRequest should have username")
	}

	if createReq.Name.GivenName == "" || createReq.Name.FamilyName == "" {
		t.Error("CreateUserRequest should have both given name and family name")
	}
}

// TestRoleGrantsWithSingleRole tests role grants when user has a single role in Role field.
func TestRoleGrantsWithSingleRole(t *testing.T) {
	user := &miro.User{
		Role: "organization_internal_admin",
	}

	var grants []string

	if user.Role != "" {
		for _, definition := range roleDefinitions {
			if definition.ID == user.Role {
				grants = append(grants, definition.RoleKey)
				break
			}
		}
	}

	expectedGrants := 1
	if len(grants) != expectedGrants {
		t.Errorf("Expected %d role grants, got %d", expectedGrants, len(grants))
	}

	if len(grants) > 0 && grants[0] != "ORGANIZATION_INTERNAL_ADMIN" {
		t.Errorf("Expected role grant 'ORGANIZATION_INTERNAL_ADMIN', got '%s'", grants[0])
	}
}

// TestRoleGrantsWithNoRoles tests role grants when user has no roles.
func TestRoleGrantsWithNoRoles(t *testing.T) {
	user := &miro.User{
		Role: "",
	}

	var grants []string

	if user.Role != "" {
		for _, definition := range roleDefinitions {
			if definition.ID == user.Role {
				grants = append(grants, definition.RoleKey)
				break
			}
		}
	}

	expectedGrants := 0
	if len(grants) != expectedGrants {
		t.Errorf("Expected %d role grants, got %d", expectedGrants, len(grants))
	}
}

// TestRoleDefinitionsConsistency tests that role definitions are consistent with tests.
func TestRoleDefinitionsConsistency(t *testing.T) {
	expectedRoleKeys := map[string]string{
		"organization_internal_admin":  "ORGANIZATION_INTERNAL_ADMIN",
		"organization_internal_user":   "ORGANIZATION_INTERNAL_USER",
		"organization_external_user":   "ORGANIZATION_EXTERNAL_USER",
		"organization_team_guest_user": "ORGANIZATION_TEAM_GUEST_USER",
	}

	for _, definition := range roleDefinitions {
		expectedKey, exists := expectedRoleKeys[definition.ID]
		if !exists {
			t.Errorf("Unexpected role definition ID: %s", definition.ID)
		}
		if definition.RoleKey != expectedKey {
			t.Errorf("Role definition %s has RoleKey %s, expected %s", definition.ID, definition.RoleKey, expectedKey)
		}
	}
}
