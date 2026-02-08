package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	MiroAccessToken = field.StringField(
		"miro-access-token",
		field.WithRequired(true),
		field.WithDescription("Miro access token. This is used to authenticate with the Miro API and sync users, teams and roles. Assign team to user and unassign user from team."),
		field.WithDisplayName("Miro Access Token"),
		field.WithIsSecret(true),
	)
	MiroScimAccessToken = field.StringField(
		"miro-scim-access-token",
		field.WithDescription("Miro SCIM access token. This is used to authenticate with the Miro SCIM API and create users. Assign role to user and revoke role from user."),
		field.WithDisplayName("Miro SCIM Access Token"),
		field.WithIsSecret(true),
	)
	BaseURLField = field.StringField(
		"base-url",
		field.WithDisplayName("Base URL"),
		field.WithDescription("Override the Miro API URL (for testing or enterprise deployments)"),
	)
	ConfigurationFields = []field.SchemaField{MiroAccessToken, MiroScimAccessToken, BaseURLField}
)

//go:generate go run ./gen
var (
	Config = field.NewConfiguration(
		ConfigurationFields,
		field.WithConnectorDisplayName("Miro"),
		field.WithHelpUrl("/docs/baton/miro"),
		field.WithIconUrl("/static/app-icons/miro.svg"),
	)
)
