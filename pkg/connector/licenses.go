package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-miro/pkg/miro"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
)

const licenseAssigned = "assigned"

type licenseDefinition struct {
	ID          string
	DisplayName string
}

var licenseDefinitions = map[string]licenseDefinition{
	"full":            {ID: "full", DisplayName: "Full License"},
	"occasional":      {ID: "occasional", DisplayName: "Occasional License"},
	"free":            {ID: "free", DisplayName: "Free License"},
	"free_restricted": {ID: "free_restricted", DisplayName: "Free Restricted License"},
	"full_trial":      {ID: "full_trial", DisplayName: "Full Trial License"},
}

type licenseBuilder struct {
	client       *miro.Client
	resourceType *v2.ResourceType
}

func (l *licenseBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return l.resourceType
}

// List returns the resources for the license builder.
func (l *licenseBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var resources []*v2.Resource
	for _, license := range licenseDefinitions {
		licenseResource, err := resource.NewResource(
			license.DisplayName,
			l.resourceType,
			license.ID,
		)
		if err != nil {
			return nil, "", nil, fmt.Errorf("failed to create license resource: %w", err)
		}
		resources = append(resources, licenseResource)
	}
	return resources, "", nil, nil
}

// Entitlements returns the entitlements for the license builder.
func (l *licenseBuilder) Entitlements(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var rv []*v2.Entitlement

	assigmentOptions := []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDescription(fmt.Sprintf("Has %s", resource.DisplayName)),
		entitlement.WithDisplayName(fmt.Sprintf("%s assigned", resource.DisplayName)),
	}

	entitlement := entitlement.NewAssignmentEntitlement(resource, licenseAssigned, assigmentOptions...)
	rv = append(rv, entitlement)

	return rv, "", nil, nil
}

// Grants returns empty grants - license grants are now emitted from user resources.
func (l *licenseBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// newLicenseBuilder creates a new license builder.
func newLicenseBuilder(client *miro.Client) *licenseBuilder {
	return &licenseBuilder{
		client:       client,
		resourceType: licenseResourceType,
	}
}
