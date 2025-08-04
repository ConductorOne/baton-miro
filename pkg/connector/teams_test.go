package connector

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/conductorone/baton-miro/pkg/miro"
	"github.com/conductorone/baton-miro/test"
)

const (
	testUserID = "user-123"
	testTeamID = "team-123"
)

// TestTeamResource tests the team resource.
func TestTeamResource(t *testing.T) {
	team := &miro.Team{
		Id:   testTeamID,
		Name: "Engineering Team",
		Type: "team",
	}

	resource, err := teamResource(team)
	if err != nil {
		t.Fatalf("teamResource() error = %v", err)
	}

	if resource.DisplayName != "Engineering Team" {
		t.Errorf("teamResource() DisplayName = %v, want %v", resource.DisplayName, "Engineering Team")
	}

	if resource.Id.Resource != testTeamID {
		t.Errorf("teamResource() Id.Resource = %v, want %v", resource.Id.Resource, testTeamID)
	}

	if resource.Id.ResourceType != teamResourceType.Id {
		t.Errorf("teamResource() Id.ResourceType = %v, want %v", resource.Id.ResourceType, teamResourceType.Id)
	}
}

// TestTeamRolesValidation tests the team roles validation.
func TestTeamRolesValidation(t *testing.T) {
	expectedRoles := []string{
		nonTeamTeamRole,
		memberTeamRole,
		adminTeamRole,
		teamGuestTeamRole,
	}

	if len(teamRoles) != len(expectedRoles) {
		t.Errorf("teamRoles length = %v, want %v", len(teamRoles), len(expectedRoles))
	}

	for i, role := range expectedRoles {
		if teamRoles[i] != role {
			t.Errorf("teamRoles[%d] = %v, want %v", i, teamRoles[i], role)
		}
	}
}

// TestTeamMockData tests the team mock data.
func TestTeamMockData(t *testing.T) {
	mockData := test.ReadFile("teams_success.json")

	var response miro.GetTeamsResponse
	err := json.Unmarshal([]byte(mockData), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal mock teams data: %v", err)
	}

	if len(response.Data) != 2 {
		t.Errorf("Expected 2 teams in mock data, got %d", len(response.Data))
	}

	if response.Data[0].Name != "Engineering Team" {
		t.Errorf("Expected first team name to be 'Engineering Team', got %s", response.Data[0].Name)
	}

	if response.Data[0].Id != testTeamID {
		t.Errorf("Expected first team ID to be %s, got %s", testTeamID, response.Data[0].Id)
	}
}

// TestTeamMembersMockData tests the team members mock data.
func TestTeamMembersMockData(t *testing.T) {
	mockData := test.ReadFile("team_members_success.json")

	var response miro.GetTeamMembersResponse
	err := json.Unmarshal([]byte(mockData), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal mock team members data: %v", err)
	}

	if len(response.Data) != 2 {
		t.Errorf("Expected 2 members in mock data, got %d", len(response.Data))
	}

	if response.Data[0].Role != "admin" {
		t.Errorf("Expected first member role to be 'admin', got %s", response.Data[0].Role)
	}

	if response.Data[0].Id != testUserID {
		t.Errorf("Expected first member ID to be %s, got %s", testUserID, response.Data[0].Id)
	}
}

// TestTeamBuilderResourceType tests the team builder resource type.
func TestTeamBuilderResourceType(t *testing.T) {
	builder := &teamBuilder{
		resourceType: teamResourceType,
	}

	result := builder.ResourceType(context.TODO())

	if result != teamResourceType {
		t.Errorf("ResourceType() = %v, want %v", result, teamResourceType)
	}
}
