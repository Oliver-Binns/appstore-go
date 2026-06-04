package users

import (
	"context"
	"fmt"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/oliver-binns/appstore-go/internal/ptr"
	"github.com/oliver-binns/appstore-go/networking"
	"github.com/oliver-binns/appstore-go/openapi"
)

func Create(c networking.HTTPClient, ctx context.Context, rawURL string, user User) (*User, error) {
	client, err := openapi.NewClientWithResponses(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	requestObject := openapi.UserInvitationCreateRequest{}
	requestObject.Data.Type = openapi.UserInvitationCreateRequestDataTypeUserInvitations
	requestObject.Data.Attributes = openapi.UserInvitationAttributes{
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		Email:               openapi_types.Email(user.Username),
		Roles:               user.Roles,
		AllAppsVisible:      ptr.PtrOrNil(user.AllAppsVisible),
		ProvisioningAllowed: ptr.PtrOrNil(user.ProvisioningAllowed),
	}
	requestObject.Data.Relationships = invitationCreateRelationships(user.VisibleAppIDs, user.AllAppsVisible)

	resp, err := client.UserInvitationsCreateInstanceWithResponse(ctx, requestObject)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 201 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	if resp.JSON201 == nil {
		return nil, fmt.Errorf("failed to decode response")
	}

	userResponse := resp.JSON201
	return &User{
		ID:                  userResponse.Data.Id,
		FirstName:           ptr.Deref(userResponse.Data.Attributes.FirstName),
		LastName:            ptr.Deref(userResponse.Data.Attributes.LastName),
		Username:            string(ptr.Deref(userResponse.Data.Attributes.Email)),
		Roles:               ptr.Deref(userResponse.Data.Attributes.Roles),
		AllAppsVisible:      ptr.Deref(userResponse.Data.Attributes.AllAppsVisible),
		ProvisioningAllowed: ptr.Deref(userResponse.Data.Attributes.ProvisioningAllowed),
		VisibleAppIDs:       user.VisibleAppIDs,
	}, nil
}
