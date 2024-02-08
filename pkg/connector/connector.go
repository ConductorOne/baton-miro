package connector

import (
	"context"
	"io"

	"github.com/conductorone/baton-miro/pkg/miro"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

type Miro struct {
	OrganizationId string

	Client *miro.Client
}

// ResourceSyncers returns a ResourceSyncer for each resource type that should be synced from the upstream service.
func (d *Miro) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		newUserBuilder(d.Client, d.OrganizationId),
		newTeamBuilder(d.Client, d.OrganizationId),
	}
}

// Asset takes an input AssetRef and attempts to fetch it using the connector's authenticated http client
// It streams a response, always starting with a metadata object, following by chunked payloads for the asset.
func (d *Miro) Asset(ctx context.Context, asset *v2.AssetRef) (string, io.ReadCloser, error) {
	return "", nil, nil
}

// Metadata returns metadata about the connector.
func (d *Miro) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "My Baton Connector",
		Description: "The template implementation of a baton connector",
	}, nil
}

// Validate is called to ensure that the connector is properly configured. It should exercise any API credentials
// to be sure that they are valid.
func (d *Miro) Validate(ctx context.Context) (annotations.Annotations, error) {
	return nil, nil
}

// New returns a new instance of the connector.
func New(ctx context.Context, accessToken string) (*Miro, error) {
	httpClient, err := uhttp.NewClient(
		ctx,
		uhttp.WithLogger(true, nil),
		uhttp.WithUserAgent("baton-miro"),
	)
	if err != nil {
		return nil, err
	}

	client := miro.New(accessToken, httpClient)

	context, _, err := client.GetContext(ctx)
	if err != nil {
		return nil, err
	}

	return &Miro{
		Client:         client,
		OrganizationId: context.Organization.Id,
	}, nil
}
