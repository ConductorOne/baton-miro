package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-miro/pkg/miro"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/helpers"
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

func (o *userBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return userResourceType
}

func userResource(ctx context.Context, user *miro.User) (*v2.Resource, error) {
	firstName, lastName := helpers.SplitFullName(user.Name)

	profile := map[string]interface{}{
		"first_name": firstName,
		"last_name":  lastName,
		"email":      user.Email,
		"login":      user.Email,
	}

	var status v2.UserTrait_Status_Status
	if user.Active {
		status = v2.UserTrait_Status_STATUS_ENABLED
	} else {
		status = v2.UserTrait_Status_STATUS_DISABLED
	}

	userTraits := []rs.UserTraitOption{
		rs.WithUserProfile(profile),
		rs.WithUserLogin(user.Email),
		rs.WithStatus(status),
	}

	resource, err := rs.NewUserResource(user.Name, userResourceType, user.Id, userTraits)
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

	response, _, err := o.client.GetOrganizationMembers(ctx, o.organizationId, cursor, resourcePageSize)
	if err != nil {
		return nil, "", nil, wrapError(err, "failed to get users")
	}

	fmt.Printf("response: %+v\n", o)
	fmt.Printf("response: %+v\n", response)

	var resources []*v2.Resource
	for _, user := range response.Data {
		resource, err := userResource(ctx, &user)
		if err != nil {
			return nil, "", nil, wrapError(err, "failed to create user resource")
		}

		resources = append(resources, resource)
	}

	if response.Cursor == "" {
		return resources, "", nil, nil
	}

	nextCursor, err := handleNextPage(bag, response.Cursor)
	if err != nil {
		return nil, "", nil, wrapError(err, "failed to create next page cursor")
	}

	return resources, nextCursor, nil, nil
}

// Entitlements always returns an empty slice for users.
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

	for _, role := range organizationRoles {
		assigmentOptions := []ent.EntitlementOption{
			ent.WithGrantableTo(userResourceType),
			ent.WithDescription(fmt.Sprintf("Has %s organization role", resource.DisplayName)),
			ent.WithDisplayName(fmt.Sprintf("%s organization role %s", resource.DisplayName, role)),
		}

		entitlement := ent.NewAssignmentEntitlement(resource, role, assigmentOptions...)
		rv = append(rv, entitlement)
	}

	return rv, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (o *userBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	licenseGrants, err := o.licenseGrants(ctx, resource)
	if err != nil {
		return nil, "", nil, err
	}

	organizationRoleGrants, err := o.organizationRoleGrants(ctx, resource)
	if err != nil {
		return nil, "", nil, err
	}

	grants := append(licenseGrants, organizationRoleGrants...)

	return grants, "", nil, nil
}

func (o *userBuilder) licenseGrants(ctx context.Context, resource *v2.Resource) ([]*v2.Grant, error) {
	user, _, err := o.client.GetOrganizationMember(ctx, o.organizationId, resource.Id.Resource)
	if err != nil {
		return nil, wrapError(err, "failed to get user")
	}

	if !contains(licenses, user.License) {
		return nil, wrapError(nil, "user does not have a valid license")
	}

	grant := grant.NewGrant(resource, user.License, resource.Id)

	return []*v2.Grant{grant}, nil
}

func (o *userBuilder) organizationRoleGrants(ctx context.Context, resource *v2.Resource) ([]*v2.Grant, error) {
	user, _, err := o.client.GetOrganizationMember(ctx, o.organizationId, resource.Id.Resource)
	if err != nil {
		return nil, wrapError(err, "failed to get user")
	}

	if !contains(organizationRoles, user.Role) {
		return nil, wrapError(nil, "user does not have a valid role")
	}

	grant := grant.NewGrant(resource, user.Role, resource.Id)

	return []*v2.Grant{grant}, nil
}

func newUserBuilder(client *miro.Client, organizationId string) *userBuilder {
	return &userBuilder{
		resourceType:   userResourceType,
		client:         client,
		organizationId: organizationId,
	}
}
