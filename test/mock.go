package test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/conductorone/baton-miro/pkg/miro"
)

// Mock constants.
const (
	MockBaseURL     = "https://mock.api.miro.com"
	MockAccessToken = "mock-access-token"
	MockOrgID       = "mock-org-id"
)

// MockClient is a mock implementation of the Miro client for testing.
type MockClient struct {
	// Team methods
	GetTeamsFunc         func(ctx context.Context, organizationId string, cursor string, limit int32) (*miro.GetTeamsResponse, *http.Response, error)
	GetTeamMembersFunc   func(ctx context.Context, organizationId string, teamId string, cursor string, limit int32) (*miro.GetTeamMembersResponse, *http.Response, error)
	InviteTeamMemberFunc func(ctx context.Context, organizationId string, teamId string, email string, role string) (*miro.InviteTeamMemberResponse, *http.Response, error)
	RemoveTeamMemberFunc func(ctx context.Context, organizationId string, teamId string, userId string) (*http.Response, error)

	// Organization methods
	GetOrganizationMemberFunc  func(ctx context.Context, organizationId string, userId string) (*miro.User, *http.Response, error)
	GetOrganizationMembersFunc func(ctx context.Context, organizationId string, cursor string, limit int32) (*miro.GetOrganizationMembersResponse, *http.Response, error)

	// User methods (SCIM)
	CreateUserFunc     func(ctx context.Context, user *miro.CreateUserRequest) (*miro.User, *http.Response, error)
	GetUserFunc        func(ctx context.Context, userId string) (*miro.ScimUser, *http.Response, error)
	UpdateUserRoleFunc func(ctx context.Context, userId string, role string) (*miro.ScimUser, *http.Response, error)

	// Context methods
	GetContextFunc func(ctx context.Context) (*miro.Context, *http.Response, error)
}

// GetTeams calls the mock method if it is defined.
func (m *MockClient) GetTeams(ctx context.Context, organizationId string, cursor string, limit int32) (*miro.GetTeamsResponse, *http.Response, error) {
	if m.GetTeamsFunc != nil {
		return m.GetTeamsFunc(ctx, organizationId, cursor, limit)
	}
	return nil, nil, nil
}

// GetTeamMembers calls the mock method if it is defined.
func (m *MockClient) GetTeamMembers(ctx context.Context, organizationId string, teamId string, cursor string, limit int32) (*miro.GetTeamMembersResponse, *http.Response, error) {
	if m.GetTeamMembersFunc != nil {
		return m.GetTeamMembersFunc(ctx, organizationId, teamId, cursor, limit)
	}
	return nil, nil, nil
}

// InviteTeamMember calls the mock method if it is defined.
func (m *MockClient) InviteTeamMember(ctx context.Context, organizationId string, teamId string, email string, role string) (*miro.InviteTeamMemberResponse, *http.Response, error) {
	if m.InviteTeamMemberFunc != nil {
		return m.InviteTeamMemberFunc(ctx, organizationId, teamId, email, role)
	}
	return nil, nil, nil
}

// RemoveTeamMember calls the mock method if it is defined.
func (m *MockClient) RemoveTeamMember(ctx context.Context, organizationId string, teamId string, userId string) (*http.Response, error) {
	if m.RemoveTeamMemberFunc != nil {
		return m.RemoveTeamMemberFunc(ctx, organizationId, teamId, userId)
	}
	return nil, nil
}

// GetOrganizationMember calls the mock method if it is defined.
func (m *MockClient) GetOrganizationMember(ctx context.Context, organizationId string, userId string) (*miro.User, *http.Response, error) {
	if m.GetOrganizationMemberFunc != nil {
		return m.GetOrganizationMemberFunc(ctx, organizationId, userId)
	}
	return nil, nil, nil
}

// GetOrganizationMembers calls the mock method if it is defined.
func (m *MockClient) GetOrganizationMembers(ctx context.Context, organizationId string, cursor string, limit int32) (*miro.GetOrganizationMembersResponse, *http.Response, error) {
	if m.GetOrganizationMembersFunc != nil {
		return m.GetOrganizationMembersFunc(ctx, organizationId, cursor, limit)
	}
	return nil, nil, nil
}

// CreateUser calls the mock method if it is defined.
func (m *MockClient) CreateUser(ctx context.Context, user *miro.CreateUserRequest) (*miro.User, *http.Response, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, user)
	}
	return nil, nil, nil
}

// GetUser calls the mock method if it is defined.
func (m *MockClient) GetUser(ctx context.Context, userId string) (*miro.ScimUser, *http.Response, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(ctx, userId)
	}
	return nil, nil, nil
}

// UpdateUserRole calls the mock method if it is defined.
func (m *MockClient) UpdateUserRole(ctx context.Context, userId string, role string) (*miro.ScimUser, *http.Response, error) {
	if m.UpdateUserRoleFunc != nil {
		return m.UpdateUserRoleFunc(ctx, userId, role)
	}
	return nil, nil, nil
}

// GetContext calls the mock method if it is defined.
func (m *MockClient) GetContext(ctx context.Context) (*miro.Context, *http.Response, error) {
	if m.GetContextFunc != nil {
		return m.GetContextFunc(ctx)
	}
	return nil, nil, nil
}

// ReadFile loads content from a JSON file from /test/mock/.
func ReadFile(fileName string) string {
	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)
	fullPath := filepath.Join(baseDir, "mock", fileName)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		panic(err)
	}
	return string(data)
}

// CreateMockResponseBody creates an io.ReadCloser with the contents of the file.
func CreateMockResponseBody(fileName string) io.ReadCloser {
	return io.NopCloser(strings.NewReader(ReadFile(fileName)))
}

// LoadMockJSON loads the content of a mock JSON file from /test/mock/ as []byte.
func LoadMockJSON(fileName string) []byte {
	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)
	fullPath := filepath.Join(baseDir, "mock", fileName)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		panic(err)
	}
	return data
}

// LoadMockStruct loads a mock JSON file and unmarshals it into the provided interface.
func LoadMockStruct(fileName string, v interface{}) {
	data := LoadMockJSON(fileName)
	if err := json.Unmarshal(data, v); err != nil {
		panic(err)
	}
}
