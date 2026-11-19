package connector

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	googlecredentialguard "github.com/conductorone/google-credential-guard"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type GoogleIdentityPlatform struct {
	client *auth.Client
}

var (
	resourceTypeUser = &v2.ResourceType{
		Id:          "user",
		DisplayName: "User",
		Traits: []v2.ResourceType_Trait{
			v2.ResourceType_TRAIT_USER,
		},
	}
)

// Metadata returns metadata about the connector.
func (gip *GoogleIdentityPlatform) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "Google Identity Platform",
		Description: "Connector syncing users from Google Identity Platform to Baton.",
	}, nil
}

// Validate hits the Google Identity Platform API to validate provided credentials.
func (gip *GoogleIdentityPlatform) Validate(ctx context.Context) (annotations.Annotations, error) {
	iter := gip.client.Users(ctx, "")
	for {
		_, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// New returns the Google Identity Platform connector.
func New(ctx context.Context, credentialsJSONFilePath string) (*GoogleIdentityPlatform, error) {
	opt, err := credentialsOption(credentialsJSONFilePath)
	if err != nil {
		return nil, err
	}
	// An explicit config prevents Firebase from reading FIREBASE_CONFIG as a file.
	app, err := firebase.NewApp(ctx, &firebase.Config{}, opt)
	if err != nil {
		return nil, errors.New("initialize Firebase app: invalid configuration")
	}

	client, err := app.Auth(ctx)
	if err != nil {
		// Firebase's signer error can print the complete private key. Do not wrap
		// or log its raw error text; keep this failure class intentionally generic.
		return nil, errors.New("initialize Firebase auth client: failed")
	}

	return &GoogleIdentityPlatform{
		client: client,
	}, nil
}

// credentialsOption reads an operator-configured local key with a strict size
// limit before handing validated service-account bytes to the Google SDK.
func credentialsOption(path string) (option.ClientOption, error) {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, errors.New("open credentials file: not found")
		}
		if errors.Is(err, fs.ErrPermission) {
			return nil, errors.New("open credentials file: permission denied")
		}
		return nil, errors.New("open credentials file: failed")
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return nil, errors.New("inspect credentials file: failed")
	}
	if !info.Mode().IsRegular() {
		return nil, errors.New("credentials file is not regular")
	}

	data, err := io.ReadAll(io.LimitReader(file, googlecredentialguard.MaxJSONBytes+1))
	if err != nil {
		return nil, errors.New("read credentials file: failed")
	}
	key, err := googlecredentialguard.ParseServiceAccountJSON(data)
	if err != nil {
		return nil, fmt.Errorf("validate credentials file: %w", err)
	}
	return key.ClientOption()
}

func (gip *GoogleIdentityPlatform) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		userBuilder(gip.client),
	}
}
