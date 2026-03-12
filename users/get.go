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

func Get(c networking.HTTPClient, ctx context.Context, rawURL string, id string) (*User, error) {
	// Parse the raw URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	parsedURL.Path = path.Join(parsedURL.Path, "users", id)

	query := parsedURL.Query()
	query.Set("include", "visibleApps")
	parsedURL.RawQuery = query.Encode()

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Do(req)

	// If the user is not found, see if the user invitation exists
	if resp.StatusCode == http.StatusNotFound {
		return getInvitations(c, ctx, rawURL, id)
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err != nil {
		return nil, err
	}

	var userResponse openapi.UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if err := resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close response body: %w", err)
	}

	u := userResponse.Data
	return &User{
		ID:                  u.Id,
		FirstName:           derefString(u.Attributes.FirstName),
		LastName:            derefString(u.Attributes.LastName),
		Username:            derefString(u.Attributes.Username),
		Roles:               derefRoles(u.Attributes.Roles),
		AllAppsVisible:      derefBool(u.Attributes.AllAppsVisible),
		ProvisioningAllowed: derefBool(u.Attributes.ProvisioningAllowed),
		HasAcceptedInvite:   true,
		VisibleAppIDs:       visibleAppIDs(u.Relationships),
	}, nil
}

func getInvitations(c networking.HTTPClient, ctx context.Context, rawURL string, id string) (*User, error) {
	parsedURL, _ := url.Parse(rawURL)
	parsedURL.Path = path.Join(parsedURL.Path, "userInvitations", id)

	query := parsedURL.Query()
	query.Set("include", "visibleApps")
	parsedURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := c.Do(req)

	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var invResponse openapi.UserInvitationResponse
	if err := json.NewDecoder(resp.Body).Decode(&invResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if err := resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close response body: %w", err)
	}

	inv := invResponse.Data
	return &User{
		ID:                  inv.Id,
		FirstName:           derefString(inv.Attributes.FirstName),
		LastName:            derefString(inv.Attributes.LastName),
		Username:            derefString(inv.Attributes.Email),
		Roles:               derefRoles(inv.Attributes.Roles),
		AllAppsVisible:      derefBool(inv.Attributes.AllAppsVisible),
		ProvisioningAllowed: derefBool(inv.Attributes.ProvisioningAllowed),
		HasAcceptedInvite:   false,
		VisibleAppIDs:       invitationVisibleAppIDs(inv.Relationships),
	}, nil
}
