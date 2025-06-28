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

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Do(req)

	// If the user is not found, see if the user invitation exists
	if resp.StatusCode == http.StatusNotFound {
		return getInvitations(c, ctx, rawURL, id)
	}

	if err != nil {
		return nil, err
	}

	userResponse := new(connectapi.Response[User])
	if err := json.NewDecoder(resp.Body).Decode(userResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if err := resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close response body: %w", err)
	}

	userResponse.Data.Data.ID = userResponse.Data.ID
	return &userResponse.Data.Data, nil
}

func getInvitations(c networking.HTTPClient, ctx context.Context, rawURL string, id string) (*User, error) {
	parsedURL, _ := url.Parse(rawURL)
	parsedURL.Path = path.Join(parsedURL.Path, "userInvitations", id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := c.Do(req)

	if err != nil {
		return nil, err
	}

	userResponse := new(connectapi.Response[userInvitation])
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
	}, nil
}
