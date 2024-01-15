package connector

import (
	"context"

	"github.com/conductorone/baton-miro/pkg/miro"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type teamBuilder struct {
	resourceType   *v2.ResourceType
	client         *miro.Client
	organizationId string
}

func (o *teamBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return teamResourceType
}

func teamResource(ctx context.Context, team *miro.Team) (*v2.Resource, error) {
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

func newTeamBuilder(client *miro.Client, organizationId string) *teamBuilder {
	return &teamBuilder{
		resourceType:   teamResourceType,
		client:         client,
		organizationId: organizationId,
	}
}

func (g *teamBuilder) List(ctx context.Context, _ *v2.ResourceId, pagination *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	bag, cursor, err := parsePageToken(pagination.Token, &v2.ResourceId{ResourceType: g.resourceType.Id})
	if err != nil {
		return nil, "", nil, wrapError(err, "failed to parse page token")
	}

	response, _, err := g.client.GetTeams(ctx, g.organizationId, cursor, resourcePageSize)
	if err != nil {
		return nil, "", nil, wrapError(err, "failed to get teams")
	}

	var resources []*v2.Resource
	for _, team := range response.Data {
		resource, err := teamResource(ctx, &team)
		if err != nil {
			return nil, "", nil, wrapError(err, "failed to create team resource")
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

func (o *teamBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (o *teamBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}
