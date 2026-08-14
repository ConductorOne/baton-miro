package main

import (
	"context"

	cfg "github.com/conductorone/baton-miro/pkg/config"
	"github.com/conductorone/baton-miro/pkg/connector"
	"github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorrunner"
)

var version = "dev"

func main() {
	ctx := context.Background()
	config.RunConnector(
		ctx,
		"baton-miro",
		version,
		cfg.Configuration,
		connector.New,
		connectorrunner.WithProvisioningEnabled(),
		connectorrunner.WithDefaultCapabilitiesConnectorBuilderV2(&connector.Connector{}),
	)
}
