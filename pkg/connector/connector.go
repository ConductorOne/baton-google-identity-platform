package connector

import (
	"context"
	"errors"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
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
func New(ctx context.Context, credentialsJSON string) (*GoogleIdentityPlatform, error) {
	opt := option.WithCredentialsFile(credentialsJSON)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return &GoogleIdentityPlatform{
		client: client,
	}, nil
}

func (gip *GoogleIdentityPlatform) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		userBuilder(gip.client),
	}
}
