package users

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/oliver-binns/appstore-go/openapi"
	"github.com/oliver-binns/googleplay-go/networking"
)

func Create(c networking.HTTPClient, ctx context.Context, rawURL string, user User) (*User, error) {
	// Parse the raw URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	parsedURL.Path = path.Join(parsedURL.Path, "userInvitations")

	requestData := openapi.UserInvitationCreateRequestData{
		Type: "userInvitations",
		Attributes: openapi.UserInvitationCreateRequestAttributes{
			FirstName:           user.FirstName,
			LastName:            user.LastName,
			Email:               user.Username,
			Roles:               user.Roles,
			AllAppsVisible:      boolPtrOrNil(user.AllAppsVisible),
			ProvisioningAllowed: boolPtrOrNil(user.ProvisioningAllowed),
		},
	}

	if linkages := user.visibleAppsLinkages(); linkages != nil {
		requestData.Relationships = &openapi.UserInvitationCreateRequestRelationships{
			VisibleApps: linkages,
		}
	}

	// Create the request body
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(openapi.UserInvitationCreateRequest{Data: requestData}); err != nil {
		return nil, fmt.Errorf("failed to encode request body: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, parsedURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusCreated {
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
		// Visible App IDs are returned from the input as these are not available in the API response:
		VisibleAppIDs: user.VisibleAppIDs,
	}, nil
}
