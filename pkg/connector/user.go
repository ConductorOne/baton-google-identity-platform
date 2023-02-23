package connector

import (
	"context"
	"errors"
	"strings"

	"firebase.google.com/go/auth"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	resource "github.com/conductorone/baton-sdk/pkg/types/resource"
	"google.golang.org/api/iterator"
)

type userResourceType struct {
	resourceType *v2.ResourceType
	client       *auth.Client
}

func (o *userResourceType) ResourceType(_ context.Context) *v2.ResourceType {
	return o.resourceType
}

// Create a new connector resource for a Google Identity Platform user.
func userResource(ctx context.Context, user *auth.ExportedUserRecord, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	names := strings.SplitN(user.DisplayName, " ", 2)
	var firstName, lastName, displayName string
	switch len(names) {
	case 1:
		firstName = names[0]
	case 2:
		firstName = names[0]
		lastName = names[1]
	}
	profile := make(map[string]interface{})
	profile["first_name"] = firstName
	profile["last_name"] = lastName
	profile["login"] = user.Email
	profile["user_id"] = user.UID

	var userStatus v2.UserTrait_Status_Status
	if user.Disabled {
		userStatus = v2.UserTrait_Status_STATUS_DISABLED
	} else {
		userStatus = v2.UserTrait_Status_STATUS_ENABLED
	}

	if user.DisplayName != "" {
		displayName = user.DisplayName
	} else {
		displayName = user.Email
	}

	userTraitOptions := []resource.UserTraitOption{resource.WithUserProfile(profile), resource.WithEmail(user.Email, true), resource.WithStatus(userStatus)}
	ret, err := resource.NewUserResource(displayName, resourceTypeUser, user.UID, userTraitOptions, resource.WithParentResourceID(parentResourceID))
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (o *userResourceType) Entitlements(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (o *userResourceType) Grants(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (o *userResourceType) List(ctx context.Context, parentResourceID *v2.ResourceId, pt *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var rv []*v2.Resource

	iter := o.client.Users(ctx, "")
	for {
		user, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, "", nil, err
		}

		ur, err := userResource(ctx, user, parentResourceID)
		if err != nil {
			return nil, "", nil, err
		}
		rv = append(rv, ur)
	}

	return rv, "", nil, nil
}

func userBuilder(client *auth.Client) *userResourceType {
	return &userResourceType{
		resourceType: resourceTypeUser,
		client:       client,
	}
}
