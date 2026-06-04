package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/oliver-binns/appstore-go/internal/ptr"
	"github.com/oliver-binns/appstore-go/networking"
	"github.com/oliver-binns/appstore-go/openapi"
)

func FindByEmail(c networking.HTTPClient, ctx context.Context, rawURL string, email string) (*User, error) {
	user, err := findActiveUserByEmail(c, ctx, rawURL, email)
	if err != nil || user != nil {
		return user, err
	}
	return findInvitationByEmail(c, ctx, rawURL, email)
}

func findActiveUserByEmail(c networking.HTTPClient, ctx context.Context, rawURL string, email string) (*User, error) {
	client, err := openapi.NewClientWithResponses(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	filters := []string{email}
	include := []openapi.UsersGetCollectionParamsInclude{openapi.UsersGetCollectionParamsIncludeVisibleApps}
	resp, err := client.UsersGetCollectionWithResponse(ctx, &openapi.UsersGetCollectionParams{
		FilterUsername: &filters,
		Include:        &include,
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("failed to decode response")
	}

	listResponse := resp.JSON200
	if len(listResponse.Data) == 0 {
		return nil, nil
	}

	return &User{
		ID:                  listResponse.Data[0].Id,
		FirstName:           ptr.Deref(listResponse.Data[0].Attributes.FirstName),
		LastName:            ptr.Deref(listResponse.Data[0].Attributes.LastName),
		Username:            ptr.Deref(listResponse.Data[0].Attributes.Username),
		Roles:               ptr.Deref(listResponse.Data[0].Attributes.Roles),
		AllAppsVisible:      ptr.Deref(listResponse.Data[0].Attributes.AllAppsVisible),
		ProvisioningAllowed: ptr.Deref(listResponse.Data[0].Attributes.ProvisioningAllowed),
		HasAcceptedInvite:   true,
		VisibleAppIDs:       visibleAppIDs(listResponse.Data[0].Relationships),
	}, nil
}

func findInvitationByEmail(c networking.HTTPClient, ctx context.Context, rawURL string, email string) (*User, error) {
	client, err := openapi.NewClientWithResponses(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	filters := []string{email}
	include := []openapi.UserInvitationsGetCollectionParamsInclude{openapi.UserInvitationsGetCollectionParamsIncludeVisibleApps}
	resp, err := client.UserInvitationsGetCollectionWithResponse(ctx, &openapi.UserInvitationsGetCollectionParams{
		FilterEmail: &filters,
		Include:     &include,
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("failed to decode response")
	}

	listResponse := resp.JSON200
	if len(listResponse.Data) == 0 {
		return nil, nil
	}

	return &User{
		ID:                  listResponse.Data[0].Id,
		FirstName:           ptr.Deref(listResponse.Data[0].Attributes.FirstName),
		LastName:            ptr.Deref(listResponse.Data[0].Attributes.LastName),
		Username:            string(ptr.Deref(listResponse.Data[0].Attributes.Email)),
		Roles:               ptr.Deref(listResponse.Data[0].Attributes.Roles),
		AllAppsVisible:      ptr.Deref(listResponse.Data[0].Attributes.AllAppsVisible),
		ProvisioningAllowed: ptr.Deref(listResponse.Data[0].Attributes.ProvisioningAllowed),
		HasAcceptedInvite:   false,
		VisibleAppIDs:       visibleAppIDs(listResponse.Data[0].Relationships),
	}, nil
}
