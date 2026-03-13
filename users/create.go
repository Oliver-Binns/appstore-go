package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/oliver-binns/appstore-go/openapi"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/oliver-binns/googleplay-go/networking"
)

func Create(c networking.HTTPClient, ctx context.Context, rawURL string, user User) (*User, error) {
	apiClient, err := openapi.NewClient(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	var req openapi.UserInvitationCreateRequest
	req.Data.Type = openapi.UserInvitationCreateRequestDataTypeUserInvitations
	req.Data.Attributes.Email = openapi_types.Email(user.Username)
	req.Data.Attributes.FirstName = user.FirstName
	req.Data.Attributes.LastName = user.LastName
	req.Data.Attributes.Roles = user.Roles
	req.Data.Attributes.AllAppsVisible = boolPtrOrNil(user.AllAppsVisible)
	req.Data.Attributes.ProvisioningAllowed = boolPtrOrNil(user.ProvisioningAllowed)

	if len(user.VisibleAppIDs) > 0 && !user.AllAppsVisible {
		apps := make([]struct {
			Id   string                                                                   `json:"id"`
			Type openapi.UserInvitationCreateRequestDataRelationshipsVisibleAppsDataType `json:"type"`
		}, len(user.VisibleAppIDs))
		for i, id := range user.VisibleAppIDs {
			apps[i].Id = id
			apps[i].Type = openapi.UserInvitationCreateRequestDataRelationshipsVisibleAppsDataTypeApps
		}
		visibleApps := struct {
			Data *[]struct {
				Id   string                                                                   `json:"id"`
				Type openapi.UserInvitationCreateRequestDataRelationshipsVisibleAppsDataType `json:"type"`
			} `json:"data,omitempty"`
		}{Data: &apps}
		rels := struct {
			VisibleApps *struct {
				Data *[]struct {
					Id   string                                                                   `json:"id"`
					Type openapi.UserInvitationCreateRequestDataRelationshipsVisibleAppsDataType `json:"type"`
				} `json:"data,omitempty"`
			} `json:"visibleApps,omitempty"`
		}{VisibleApps: &visibleApps}
		req.Data.Relationships = &rels
	}

	resp, err := apiClient.UserInvitationsCreateInstance(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var invResponse openapi.UserInvitationResponse
	if err := json.NewDecoder(resp.Body).Decode(&invResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	inv := invResponse.Data
	result := &User{
		ID:            inv.Id,
		VisibleAppIDs: user.VisibleAppIDs,
	}
	if inv.Attributes != nil {
		result.FirstName = derefString(inv.Attributes.FirstName)
		result.LastName = derefString(inv.Attributes.LastName)
		result.Username = derefEmail(inv.Attributes.Email)
		result.Roles = derefRoles(inv.Attributes.Roles)
		result.AllAppsVisible = derefBool(inv.Attributes.AllAppsVisible)
		result.ProvisioningAllowed = derefBool(inv.Attributes.ProvisioningAllowed)
	}
	return result, nil
}
