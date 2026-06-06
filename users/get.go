package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/oliver-binns/appstore-go/internal/ptr"
	"github.com/oliver-binns/appstore-go/openapi"
	"github.com/oliver-binns/appstore-go/networking"
)

func Get(c networking.HTTPClient, ctx context.Context, rawURL string, id string) (*User, error) {
	if id == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

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

	userResponse := new(openapi.UserResponse)
	if err := json.NewDecoder(resp.Body).Decode(userResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if err := resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close response body: %w", err)
	}

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

	userResponse := new(openapi.UserInvitationResponse)
	if err := json.NewDecoder(resp.Body).Decode(userResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if err := resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close response body: %w", err)
	}

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
