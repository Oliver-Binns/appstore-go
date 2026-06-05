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

	visibleAppIDs := visibleAppIDs(listResponse.Data[0].Relationships)

	return &User{
		ID:                  listResponse.Data[0].Id,
		FirstName:           ptr.Deref(listResponse.Data[0].Attributes.FirstName),
		LastName:            ptr.Deref(listResponse.Data[0].Attributes.LastName),
		Username:            ptr.Deref(listResponse.Data[0].Attributes.Username),
		Roles:               ptr.Deref(listResponse.Data[0].Attributes.Roles),
		AllAppsVisible:      ptr.Deref(listResponse.Data[0].Attributes.AllAppsVisible),
		ProvisioningAllowed: ptr.Deref(listResponse.Data[0].Attributes.ProvisioningAllowed),
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

	visibleAppIDs := visibleAppIDs(listResponse.Data[0].Relationships)

	return &User{
		ID:                  listResponse.Data[0].Id,
		FirstName:           ptr.Deref(listResponse.Data[0].Attributes.FirstName),
		LastName:            ptr.Deref(listResponse.Data[0].Attributes.LastName),
		Username:            string(ptr.Deref(listResponse.Data[0].Attributes.Email)),
		Roles:               ptr.Deref(listResponse.Data[0].Attributes.Roles),
		AllAppsVisible:      ptr.Deref(listResponse.Data[0].Attributes.AllAppsVisible),
		ProvisioningAllowed: ptr.Deref(listResponse.Data[0].Attributes.ProvisioningAllowed),
		HasAcceptedInvite:   false,
		VisibleAppIDs:       visibleAppIDs,
	}, nil
}
