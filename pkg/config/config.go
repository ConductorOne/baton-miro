package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	MiroAccessToken = field.StringField(
		"miro-access-token",
		field.WithRequired(true),
		field.WithDescription("Miro access token. This is used to authenticate with the Miro API and sync users, teams and roles. Assign team to user and desasgine user from team."),
		field.WithDisplayName("Miro Access Token"),
	)
	MiroScimAccessToken = field.StringField(
		"miro-scim-access-token",
		field.WithDescription("Miro SCIM access token. This is used to authenticate with the Miro SCIM API and create users. Assign role to user and revoke role from user."),
		field.WithDisplayName("Miro SCIM Access Token"),
	)
	ConfigurationFields = []field.SchemaField{MiroAccessToken, MiroScimAccessToken}
)

var (
	Config = field.NewConfiguration(
		ConfigurationFields,
		field.WithConnectorDisplayName("Miro"),
		field.WithHelpUrl("/docs/baton/miro"),
		field.WithIconUrl("/static/app-icons/miro.svg"),
	)
)
