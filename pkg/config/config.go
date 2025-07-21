package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	MiroAccessToken = field.StringField(
		"miro-access-token",
		field.WithRequired(true),
		field.WithDescription("Miro access token"),
		field.WithDisplayName("Miro Access Token"),
	)
	ConfigurationFields = []field.SchemaField{
		MiroAccessToken,
	}
	FieldRelationships = []field.SchemaFieldRelationship{
		field.FieldsRequiredTogether(MiroAccessToken),
	}
)

var (
	Config = field.NewConfiguration(
		ConfigurationFields,
		field.WithConnectorDisplayName("Miro"),
		field.WithHelpUrl("/docs/baton/miro"),
		field.WithIconUrl("/static/app-icons/miro.svg"),
	)
)
