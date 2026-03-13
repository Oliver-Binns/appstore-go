package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/oliver-binns/appstore-go/openapi"
	"github.com/oliver-binns/googleplay-go/networking"
)

func Get(c networking.HTTPClient, ctx context.Context, rawURL string, id string) (*User, error) {
	apiClient, err := openapi.NewClient(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	include := []openapi.UsersGetInstanceParamsInclude{openapi.UsersGetInstanceParamsIncludeVisibleApps}
	resp, err := apiClient.UsersGetInstance(ctx, id, &openapi.UsersGetInstanceParams{Include: &include})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return getInvitations(c, ctx, rawURL, id)
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var userResponse openapi.UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	u := userResponse.Data
	user := &User{ID: u.Id, HasAcceptedInvite: true}
	if u.Attributes != nil {
		user.FirstName = derefString(u.Attributes.FirstName)
		user.LastName = derefString(u.Attributes.LastName)
		user.Username = derefString(u.Attributes.Username)
		user.Roles = derefRoles(u.Attributes.Roles)
		user.AllAppsVisible = derefBool(u.Attributes.AllAppsVisible)
		user.ProvisioningAllowed = derefBool(u.Attributes.ProvisioningAllowed)
	}
	if u.Relationships != nil && u.Relationships.VisibleApps != nil && u.Relationships.VisibleApps.Data != nil {
		for _, app := range *u.Relationships.VisibleApps.Data {
			user.VisibleAppIDs = append(user.VisibleAppIDs, app.Id)
		}
	}
	return user, nil
}

func getInvitations(c networking.HTTPClient, ctx context.Context, rawURL string, id string) (*User, error) {
	apiClient, err := openapi.NewClient(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	include := []openapi.UserInvitationsGetInstanceParamsInclude{openapi.UserInvitationsGetInstanceParamsIncludeVisibleApps}
	resp, err := apiClient.UserInvitationsGetInstance(ctx, id, &openapi.UserInvitationsGetInstanceParams{Include: &include})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var invResponse openapi.UserInvitationResponse
	if err := json.NewDecoder(resp.Body).Decode(&invResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	inv := invResponse.Data
	user := &User{ID: inv.Id, HasAcceptedInvite: false}
	if inv.Attributes != nil {
		user.FirstName = derefString(inv.Attributes.FirstName)
		user.LastName = derefString(inv.Attributes.LastName)
		user.Username = derefEmail(inv.Attributes.Email)
		user.Roles = derefRoles(inv.Attributes.Roles)
		user.AllAppsVisible = derefBool(inv.Attributes.AllAppsVisible)
		user.ProvisioningAllowed = derefBool(inv.Attributes.ProvisioningAllowed)
	}
	if inv.Relationships != nil && inv.Relationships.VisibleApps != nil && inv.Relationships.VisibleApps.Data != nil {
		for _, app := range *inv.Relationships.VisibleApps.Data {
			user.VisibleAppIDs = append(user.VisibleAppIDs, app.Id)
		}
	}
	return user, nil
}
