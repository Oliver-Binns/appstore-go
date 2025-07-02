package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/oliver-binns/appstore-go/connectapi"
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

	userResponse := new(connectapi.Response[User, *userRelationships])
	if err := json.NewDecoder(resp.Body).Decode(userResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if err := resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close response body: %w", err)
	}

	userResponse.Data.Data.ID = userResponse.Data.ID
	userResponse.Data.Data.HasAcceptedInvite = true
	userResponse.Data.Data.VisibleAppIDs = userResponse.Data.Relationships.ids()
	return &userResponse.Data.Data, nil
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

	userResponse := new(connectapi.Response[userInvitation, *userRelationships])
	if err := json.NewDecoder(resp.Body).Decode(userResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if err := resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close response body: %w", err)
	}

	return &User{
		ID:                  userResponse.Data.ID,
		FirstName:           userResponse.Data.Data.FirstName,
		LastName:            userResponse.Data.Data.LastName,
		Username:            userResponse.Data.Data.Email,
		Roles:               userResponse.Data.Data.Roles,
		AllAppsVisible:      userResponse.Data.Data.AllAppsVisible,
		ProvisioningAllowed: userResponse.Data.Data.ProvisioningAllowed,
		HasAcceptedInvite:   false,
		VisibleAppIDs:       userResponse.Data.Relationships.ids(),
	}, nil
}
