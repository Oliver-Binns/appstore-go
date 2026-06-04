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
	if id == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	// Parse the raw URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	parsedURL.Path = path.Join(parsedURL.Path, "users", id)

	// Create the request body
	body := bytes.NewBuffer(nil)
	requestObject := openapi.UserUpdateRequest{}
	requestObject.Data.Id = id
	requestObject.Data.Type = openapi.UserUpdateRequestDataTypeUsers
	requestObject.Data.Attributes = &struct {
		AllAppsVisible      *bool               `json:"allAppsVisible,omitempty"`
		ProvisioningAllowed *bool               `json:"provisioningAllowed,omitempty"`
		Roles               *[]openapi.UserRole `json:"roles,omitempty"`
	}{
		AllAppsVisible:      boolPtrOrNil(user.AllAppsVisible),
		ProvisioningAllowed: boolPtrOrNil(user.ProvisioningAllowed),
		Roles:               rolesOrNil(user.Roles),
	}
	if len(user.VisibleAppIDs) != 0 || !user.AllAppsVisible {
		appReferences := make([]struct {
			Id   string                                                        `json:"id"`
			Type openapi.UserUpdateRequestDataRelationshipsVisibleAppsDataType `json:"type"`
		}, len(user.VisibleAppIDs))
		for index, id := range user.VisibleAppIDs {
			appReferences[index].Id = id
			appReferences[index].Type = openapi.Apps
		}
		requestObject.Data.Relationships = &struct {
			VisibleApps *struct {
				Data *[]struct {
					Id   string                                                        `json:"id"`
					Type openapi.UserUpdateRequestDataRelationshipsVisibleAppsDataType `json:"type"`
				} `json:"data,omitempty"`
			} `json:"visibleApps,omitempty"`
		}{
			VisibleApps: &struct {
				Data *[]struct {
					Id   string                                                        `json:"id"`
					Type openapi.UserUpdateRequestDataRelationshipsVisibleAppsDataType `json:"type"`
				} `json:"data,omitempty"`
			}{
				Data: &appReferences,
			},
		}
	}
	err = json.NewEncoder(body).Encode(requestObject)
	if err != nil {
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

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	userResponse := new(openapi.UserResponse)
	if err := json.NewDecoder(resp.Body).Decode(userResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if err := resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close response body: %w", err)
	}

	var firstName, lastName, username string
	var roles []openapi.UserRole
	var allAppsVisible, provisioningAllowed bool
	if userResponse.Data.Attributes != nil {
		firstName = derefString(userResponse.Data.Attributes.FirstName)
		lastName = derefString(userResponse.Data.Attributes.LastName)
		username = derefString(userResponse.Data.Attributes.Username)
		roles = derefRoles(userResponse.Data.Attributes.Roles)
		allAppsVisible = derefBool(userResponse.Data.Attributes.AllAppsVisible)
		provisioningAllowed = derefBool(userResponse.Data.Attributes.ProvisioningAllowed)
	}

	visibleAppIDs := []string{}
	visibleAppIDs = append(visibleAppIDs, user.VisibleAppIDs...)

	return &User{
		ID:                  userResponse.Data.Id,
		FirstName:           firstName,
		LastName:            lastName,
		Username:            username,
		Roles:               roles,
		AllAppsVisible:      allAppsVisible,
		ProvisioningAllowed: provisioningAllowed,
		// Visible App IDs are returned from the input as these are not available in the API response:
		VisibleAppIDs: visibleAppIDs,
	}, nil
}
