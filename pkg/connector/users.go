package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-miro/pkg/miro"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	grant "github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type userBuilder struct {
	resourceType   *v2.ResourceType
	client         *miro.Client
	organizationId string
}

func (b *userBuilder) CreateAccountCapabilityDetails(_ context.Context) (*v2.CredentialDetailsAccountProvisioning, annotations.Annotations, error) {
	return &v2.CredentialDetailsAccountProvisioning{
		SupportedCredentialOptions: []v2.CapabilityDetailCredentialOption{
			v2.CapabilityDetailCredentialOption_CAPABILITY_DETAIL_CREDENTIAL_OPTION_NO_PASSWORD,
		},
		PreferredCredentialOption: v2.CapabilityDetailCredentialOption_CAPABILITY_DETAIL_CREDENTIAL_OPTION_NO_PASSWORD,
	}, nil, nil
}

func (o *userBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return userResourceType
}

func userResource(user *miro.User) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"email": user.Email,
		"login": user.Email,
	}

	var status v2.UserTrait_Status_Status
	if user.Active {
		status = v2.UserTrait_Status_STATUS_ENABLED
	} else {
		status = v2.UserTrait_Status_STATUS_DISABLED
	}

	lastLogin, err := parseTime(user.LastActivityAt)
	if err != nil {
		return nil, wrapError(err, "failed to parse last login time")
	}

	userTraits := []rs.UserTraitOption{
		rs.WithUserProfile(profile),
		rs.WithUserLogin(user.Email),
		rs.WithStatus(status),
	}
	if lastLogin != nil {
		userTraits = append(userTraits, rs.WithLastLogin(*lastLogin))
	}

	resource, err := rs.NewUserResource(user.Email, userResourceType, user.Id, userTraits)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (o *userBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	bag, cursor, err := parsePageToken(pToken.Token, &v2.ResourceId{ResourceType: o.resourceType.Id})
	if err != nil {
		return nil, "", nil, wrapError(err, "failed to parse page token")
	}

	response, annos, err := o.client.GetOrganizationMembers(ctx, o.organizationId, cursor, resourcePageSize)
	if err != nil {
		return nil, "", annos, wrapError(err, "failed to get users")
	}

	var resources []*v2.Resource
	for _, user := range response.Data {
		resource, err := userResource(&user)
		if err != nil {
			return nil, "", annos, wrapError(err, "failed to create user resource")
		}

		resources = append(resources, resource)
	}

	if response.Cursor == "" {
		return resources, "", annos, nil
	}

	nextCursor, err := handleNextPage(bag, response.Cursor)
	if err != nil {
		return nil, "", annos, wrapError(err, "failed to create next page cursor")
	}

	return resources, nextCursor, nil, nil
}

// Entitlements returns license entitlements for users.
func (o *userBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var rv []*v2.Entitlement

	for _, license := range licenses {
		assigmentOptions := []ent.EntitlementOption{
			ent.WithGrantableTo(userResourceType),
			ent.WithDescription(fmt.Sprintf("Has %s license", resource.DisplayName)),
			ent.WithDisplayName(fmt.Sprintf("%s is %s license owner", resource.DisplayName, license)),
		}

		entitlement := ent.NewAssignmentEntitlement(resource, license, assigmentOptions...)
		rv = append(rv, entitlement)
	}

	return rv, "", nil, nil
}

// Grants returns license grants and role grants for users.
func (o *userBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	var allGrants []*v2.Grant

	licenseGrants, annos, err := o.licenseGrants(ctx, resource)
	if err != nil {
		return nil, "", annos, err
	}
	allGrants = append(allGrants, licenseGrants...)

	roleGrants, annos, err := o.roleGrants(ctx, resource)
	if err != nil {
		return nil, "", annos, err
	}
	allGrants = append(allGrants, roleGrants...)

	return allGrants, "", annos, nil
}

// CreateAccount creates a new user in Miro using the SCIM API.
func (o *userBuilder) CreateAccount(
	ctx context.Context,
	accountInfo *v2.AccountInfo,
	_ *v2.CredentialOptions,
) (
	connectorbuilder.CreateAccountResponse,
	[]*v2.PlaintextData,
	annotations.Annotations,
	error,
) {
	profile := accountInfo.GetProfile().AsMap()
	requiredFields := map[string]string{
		"first_name": "first_name is required",
		"last_name":  "last_name is required",
		"email":      "email is required",
	}

	for field, errMsg := range requiredFields {
		if val, ok := profile[field].(string); !ok || val == "" {
			return nil, nil, nil, fmt.Errorf("%s", errMsg)
		}
	}

	newUser, annos, err := o.client.CreateUser(ctx, profile["email"].(string), profile["first_name"].(string), profile["last_name"].(string))
	if err != nil {
		return nil, nil, annos, wrapError(err, "failed to create miro user")
	}

	resource, err := userResource(newUser)
	if err != nil {
		return nil, nil, annos, wrapError(err, "failed to create user resource from miro user")
	}

	return &v2.CreateAccountResponse_SuccessResult{
		Resource: resource,
	}, nil, annos, nil
}

// licenseGrants returns a grant for the user's license.
func (o *userBuilder) licenseGrants(ctx context.Context, resource *v2.Resource) ([]*v2.Grant, annotations.Annotations, error) {
	user, annos, err := o.client.GetOrganizationMember(ctx, o.organizationId, resource.Id.Resource)
	if err != nil {
		return nil, annos, wrapError(err, "failed to get user")
	}

	if user.License == "" {
		return nil, annos, wrapError(nil, "user does not have a valid license")
	}

	if !contains(licenses, user.License) {
		return nil, annos, wrapError(nil, "user does not have a valid license")
	}
	grant := grant.NewGrant(resource, user.License, resource.Id)

	return []*v2.Grant{grant}, annos, nil
}

// roleGrants returns grants for the user's roles.
func (o *userBuilder) roleGrants(ctx context.Context, resource *v2.Resource) ([]*v2.Grant, annotations.Annotations, error) {
	user, annos, err := o.client.GetOrganizationMember(ctx, o.organizationId, resource.Id.Resource)
	if err != nil {
		return nil, annos, wrapError(err, "failed to get user")
	}

	var grants []*v2.Grant
	if user.Role != "" {
		if definition, exists := roleDefinitions[user.Role]; exists {
			roleResource := &v2.ResourceId{
				ResourceType: roleResourceType.Id,
				Resource:     definition.ID,
			}
			g := grant.NewGrant(&v2.Resource{Id: roleResource}, assignedRole, resource.Id)
			grants = append(grants, g)
		}
	}
	return grants, annos, nil
}

func newUserBuilder(client *miro.Client, organizationId string) *userBuilder {
	return &userBuilder{
		resourceType:   userResourceType,
		client:         client,
		organizationId: organizationId,
	}
}
