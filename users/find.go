package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

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
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	parsedURL.Path = path.Join(parsedURL.Path, "users")

	query := parsedURL.Query()
	query.Set("filter[username]", email)
	query.Set("include", "visibleApps")
	parsedURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	listResponse := new(openapi.UsersResponse)
	if err := json.NewDecoder(resp.Body).Decode(listResponse); err != nil {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if err := resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close response body: %w", err)
	}

	if len(listResponse.Data) == 0 {
		return nil, nil
	}

	visibleAppIDs := []string{}
	if listResponse.Data[0].Relationships != nil && listResponse.Data[0].Relationships.VisibleApps != nil && listResponse.Data[0].Relationships.VisibleApps.Data != nil {
		for _, app := range *listResponse.Data[0].Relationships.VisibleApps.Data {
			visibleAppIDs = append(visibleAppIDs, app.Id)
		}
	}

	var firstName, lastName, username string
	var roles []openapi.UserRole
	var allAppsVisible, provisioningAllowed bool
	if listResponse.Data[0].Attributes != nil {
		firstName = derefString(listResponse.Data[0].Attributes.FirstName)
		lastName = derefString(listResponse.Data[0].Attributes.LastName)
		username = derefString(listResponse.Data[0].Attributes.Username)
		roles = derefRoles(listResponse.Data[0].Attributes.Roles)
		allAppsVisible = derefBool(listResponse.Data[0].Attributes.AllAppsVisible)
		provisioningAllowed = derefBool(listResponse.Data[0].Attributes.ProvisioningAllowed)
	}

	return &User{
		ID:                  listResponse.Data[0].Id,
		FirstName:           firstName,
		LastName:            lastName,
		Username:            username,
		Roles:               roles,
		AllAppsVisible:      allAppsVisible,
		ProvisioningAllowed: provisioningAllowed,
		HasAcceptedInvite:   true,
		VisibleAppIDs:       visibleAppIDs,
	}, nil
}

func findInvitationByEmail(c networking.HTTPClient, ctx context.Context, rawURL string, email string) (*User, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	parsedURL.Path = path.Join(parsedURL.Path, "userInvitations")

	query := parsedURL.Query()
	query.Set("filter[email]", email)
	query.Set("include", "visibleApps")
	parsedURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	listResponse := new(openapi.UserInvitationsResponse)
	if err := json.NewDecoder(resp.Body).Decode(listResponse); err != nil {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if err := resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close response body: %w", err)
	}

	if len(listResponse.Data) == 0 {
		return nil, nil
	}

	visibleAppIDs := []string{}
	if listResponse.Data[0].Relationships != nil && listResponse.Data[0].Relationships.VisibleApps != nil && listResponse.Data[0].Relationships.VisibleApps.Data != nil {
		for _, app := range *listResponse.Data[0].Relationships.VisibleApps.Data {
			visibleAppIDs = append(visibleAppIDs, app.Id)
		}
	}

	var firstName, lastName, username string
	var roles []openapi.UserRole
	var allAppsVisible, provisioningAllowed bool
	if listResponse.Data[0].Attributes != nil {
		firstName = derefString(listResponse.Data[0].Attributes.FirstName)
		lastName = derefString(listResponse.Data[0].Attributes.LastName)
		username = derefEmail(listResponse.Data[0].Attributes.Email)
		roles = derefRoles(listResponse.Data[0].Attributes.Roles)
		allAppsVisible = derefBool(listResponse.Data[0].Attributes.AllAppsVisible)
		provisioningAllowed = derefBool(listResponse.Data[0].Attributes.ProvisioningAllowed)
	}

	return &User{
		ID:                  listResponse.Data[0].Id,
		FirstName:           firstName,
		LastName:            lastName,
		Username:            username,
		Roles:               roles,
		AllAppsVisible:      allAppsVisible,
		ProvisioningAllowed: provisioningAllowed,
		HasAcceptedInvite:   false,
		VisibleAppIDs:       visibleAppIDs,
	}, nil
}
