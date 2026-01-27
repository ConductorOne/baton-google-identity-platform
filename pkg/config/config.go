package config

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-sdk/pkg/field"
)

var CredentialsJSONFilePath = field.StringField(
	"credentials-json-file-path",
	field.WithRequired(true),
	field.WithDescription("JSON credentials file name for the Google identity platform account."),
)

var ConfigurationFields = []field.SchemaField{
	CredentialsJSONFilePath,
}

var Config = field.NewConfiguration(ConfigurationFields)

func ValidateConfig(ctx context.Context, cfg *GoogleIdentityPlatform) error {
	if cfg.CredentialsJSONFilePath == "" {
		return fmt.Errorf("credentials json file path is missing")
	}
	return nil
}
