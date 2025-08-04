package config

import (
	"testing"

	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/stretchr/testify/assert"
)

// TestValidateConfig tests the validation of the Miro configuration.
func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Miro
		wantErr bool
	}{
		{
			name: "valid config with access token",
			config: &Miro{
				AccessToken: "test-access-token",
			},
			wantErr: false,
		},
		{
			name: "valid config with access token and scim access token",
			config: &Miro{
				AccessToken:     "test-access-token",
				ScimAccessToken: "test-scim-access-token",
			},
			wantErr: false,
		},
		{
			name:    "invalid config - missing access token",
			config:  &Miro{},
			wantErr: true,
		},
		{
			name: "invalid config - empty access token",
			config: &Miro{
				AccessToken: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := field.Validate(Config, tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
