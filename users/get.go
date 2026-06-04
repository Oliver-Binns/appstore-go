package users

import (
	"context"
	"fmt"

	"github.com/oliver-binns/appstore-go/internal/ptr"
	"github.com/oliver-binns/appstore-go/networking"
	"github.com/oliver-binns/appstore-go/openapi"
)

func Get(c networking.HTTPClient, ctx context.Context, rawURL string, id string) (*User, error) {
	if id == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	client, err := openapi.NewClientWithResponses(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	include := []openapi.UsersGetInstanceParamsInclude{openapi.UsersGetInstanceParamsIncludeVisibleApps}
	resp, err := client.UsersGetInstanceWithResponse(ctx, id, &openapi.UsersGetInstanceParams{Include: &include})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == 404 {
		return getInvitations(c, ctx, rawURL, id)
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
		HasAcceptedInvite:   true,
		VisibleAppIDs:       visibleAppIDs(userResponse.Data.Relationships),
	}, nil
}

func getInvitations(c networking.HTTPClient, ctx context.Context, rawURL string, id string) (*User, error) {
	client, err := openapi.NewClientWithResponses(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	include := []openapi.UserInvitationsGetInstanceParamsInclude{openapi.UserInvitationsGetInstanceParamsIncludeVisibleApps}
	resp, err := client.UserInvitationsGetInstanceWithResponse(ctx, id, &openapi.UserInvitationsGetInstanceParams{Include: &include})
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
		Username:            string(ptr.Deref(userResponse.Data.Attributes.Email)),
		Roles:               ptr.Deref(userResponse.Data.Attributes.Roles),
		AllAppsVisible:      ptr.Deref(userResponse.Data.Attributes.AllAppsVisible),
		ProvisioningAllowed: ptr.Deref(userResponse.Data.Attributes.ProvisioningAllowed),
		HasAcceptedInvite:   false,
		VisibleAppIDs:       visibleAppIDs(userResponse.Data.Relationships),
	}, nil
}
