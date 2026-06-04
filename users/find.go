package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/oliver-binns/appstore-go/openapi"
	"github.com/oliver-binns/googleplay-go/networking"
)

func FindByEmail(c networking.HTTPClient, ctx context.Context, rawURL string, email string) (*User, error) {
	user, err := findActiveUserByEmail(c, ctx, rawURL, email)
	if err != nil || user != nil {
		return user, err
	}
	return findInvitationByEmail(c, ctx, rawURL, email)
}

func findActiveUserByEmail(c networking.HTTPClient, ctx context.Context, rawURL string, email string) (*User, error) {
	apiClient, err := openapi.NewClient(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	include := []openapi.UsersGetCollectionParamsInclude{openapi.UsersGetCollectionParamsIncludeVisibleApps}
	resp, err := apiClient.UsersGetCollection(ctx, &openapi.UsersGetCollectionParams{
		FilterUsername: &[]string{email},
		Include:        &include,
	})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var listResponse openapi.UsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResponse); err != nil {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	if len(listResponse.Data) == 0 {
		return nil, nil
	}

	u := listResponse.Data[0]
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

func findInvitationByEmail(c networking.HTTPClient, ctx context.Context, rawURL string, email string) (*User, error) {
	apiClient, err := openapi.NewClient(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	include := []openapi.UserInvitationsGetCollectionParamsInclude{openapi.UserInvitationsGetCollectionParamsIncludeVisibleApps}
	resp, err := apiClient.UserInvitationsGetCollection(ctx, &openapi.UserInvitationsGetCollectionParams{
		FilterEmail: &[]string{email},
		Include:     &include,
	})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var listResponse openapi.UserInvitationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResponse); err != nil {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	if len(listResponse.Data) == 0 {
		return nil, nil
	}

	inv := listResponse.Data[0]
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
