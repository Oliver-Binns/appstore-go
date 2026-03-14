package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/oliver-binns/appstore-go/openapi"
	"github.com/oliver-binns/googleplay-go/networking"
)

func Modify(c networking.HTTPClient, ctx context.Context, rawURL string, id string, user User) (*User, error) {
	apiClient, err := openapi.NewClient(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	var req openapi.UserUpdateRequest
	req.Data.Id = id
	req.Data.Type = openapi.UserUpdateRequestDataTypeUsers
	req.Data.Attributes = &struct {
		AllAppsVisible      *bool       `json:"allAppsVisible,omitempty"`
		ProvisioningAllowed *bool       `json:"provisioningAllowed,omitempty"`
		Roles               *[]UserRole `json:"roles,omitempty"`
	}{
		Roles:               rolesOrNil(user.Roles),
		AllAppsVisible:      boolPtrOrNil(user.AllAppsVisible),
		ProvisioningAllowed: boolPtrOrNil(user.ProvisioningAllowed),
	}

	if len(user.VisibleAppIDs) > 0 && !user.AllAppsVisible {
		apps := make([]struct {
			Id   string                                                `json:"id"`
			Type openapi.UserUpdateRequestDataRelationshipsVisibleAppsDataType `json:"type"`
		}, len(user.VisibleAppIDs))
		for i, appID := range user.VisibleAppIDs {
			apps[i].Id = appID
			apps[i].Type = openapi.Apps
		}
		visibleApps := struct {
			Data *[]struct {
				Id   string                                                `json:"id"`
				Type openapi.UserUpdateRequestDataRelationshipsVisibleAppsDataType `json:"type"`
			} `json:"data,omitempty"`
		}{Data: &apps}
		rels := struct {
			VisibleApps *struct {
				Data *[]struct {
					Id   string                                                `json:"id"`
					Type openapi.UserUpdateRequestDataRelationshipsVisibleAppsDataType `json:"type"`
				} `json:"data,omitempty"`
			} `json:"visibleApps,omitempty"`
		}{VisibleApps: &visibleApps}
		req.Data.Relationships = &rels
	}

	resp, err := apiClient.UsersUpdateInstance(ctx, id, req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var userResponse openapi.UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	u := userResponse.Data
	result := &User{
		ID:            u.Id,
		// Visible App IDs are returned from the input as these are not available in the API response:
		VisibleAppIDs: user.VisibleAppIDs,
	}
	if u.Attributes != nil {
		result.FirstName = derefString(u.Attributes.FirstName)
		result.LastName = derefString(u.Attributes.LastName)
		result.Username = derefString(u.Attributes.Username)
		result.Roles = derefRoles(u.Attributes.Roles)
		result.AllAppsVisible = derefBool(u.Attributes.AllAppsVisible)
		result.ProvisioningAllowed = derefBool(u.Attributes.ProvisioningAllowed)
	}
	return result, nil
}
