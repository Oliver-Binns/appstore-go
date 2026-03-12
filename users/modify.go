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

func Modify(c networking.HTTPClient, ctx context.Context, rawURL string, id string, user User) (*User, error) {
	// Parse the raw URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	parsedURL.Path = path.Join(parsedURL.Path, "users", id)

	requestData := openapi.UserUpdateRequestData{
		Id:   id,
		Type: "users",
		Attributes: openapi.UserUpdateRequestAttributes{
			Roles:               rolesOrNil(user.Roles),
			AllAppsVisible:      boolPtrOrNil(user.AllAppsVisible),
			ProvisioningAllowed: boolPtrOrNil(user.ProvisioningAllowed),
		},
	}

	if linkages := user.visibleAppsLinkages(); linkages != nil {
		requestData.Relationships = &openapi.UserUpdateRequestRelationships{
			VisibleApps: linkages,
		}
	}

	// Create the request body
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(openapi.UserUpdateRequest{Data: requestData}); err != nil {
		return nil, fmt.Errorf("failed to encode request body: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, parsedURL.String(), body)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Do(req)

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
		// Visible App IDs are returned from the input as these are not available in the API response:
		VisibleAppIDs: user.VisibleAppIDs,
	}, nil
}
