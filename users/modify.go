package users

import (
	"context"
	"fmt"

	"github.com/oliver-binns/appstore-go/internal/ptr"
	"github.com/oliver-binns/appstore-go/networking"
	"github.com/oliver-binns/appstore-go/openapi"
)

func Modify(c networking.HTTPClient, ctx context.Context, rawURL string, id string, user User) (*User, error) {
	if id == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	client, err := openapi.NewClientWithResponses(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	requestObject := openapi.UserUpdateRequest{}
	requestObject.Data.Id = id
	requestObject.Data.Type = openapi.UserUpdateRequestDataTypeUsers
	requestObject.Data.Attributes = &openapi.UserUpdateAttributes{
		AllAppsVisible:      ptr.PtrOrNil(user.AllAppsVisible),
		ProvisioningAllowed: ptr.PtrOrNil(user.ProvisioningAllowed),
		Roles:               ptr.SlicePtrOrNil(user.Roles),
	}
	requestObject.Data.Relationships = userUpdateRelationships(user.VisibleAppIDs, user.AllAppsVisible)

	resp, err := client.UsersUpdateInstanceWithResponse(ctx, id, requestObject)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("failed to decode response")
	}

	userResponse := resp.JSON200
	return &User{
		ID:                  userResponse.Data.Id,
		FirstName:           ptr.Deref(userResponse.Data.Attributes.FirstName),
		LastName:            ptr.Deref(userResponse.Data.Attributes.LastName),
		Username:            ptr.Deref(userResponse.Data.Attributes.Username),
		Roles:               ptr.Deref(userResponse.Data.Attributes.Roles),
		AllAppsVisible:      ptr.Deref(userResponse.Data.Attributes.AllAppsVisible),
		ProvisioningAllowed: ptr.Deref(userResponse.Data.Attributes.ProvisioningAllowed),
		VisibleAppIDs:       user.VisibleAppIDs,
	}, nil
}
