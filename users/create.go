package users

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	openapi_types "github.com/oapi-codegen/runtime/types"
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

	// Create the request body
	body := bytes.NewBuffer(nil)
	requestObject := openapi.UserInvitationCreateRequest{}
	requestObject.Data.Type = openapi.UserInvitationCreateRequestDataTypeUserInvitations
	requestObject.Data.Attributes.FirstName = user.FirstName
	requestObject.Data.Attributes.LastName = user.LastName
	requestObject.Data.Attributes.Email = openapi_types.Email(user.Username)
	requestObject.Data.Attributes.Roles = user.Roles
	requestObject.Data.Attributes.AllAppsVisible = boolPtrOrNil(user.AllAppsVisible)
	requestObject.Data.Attributes.ProvisioningAllowed = boolPtrOrNil(user.ProvisioningAllowed)
	if len(user.VisibleAppIDs) != 0 || !user.AllAppsVisible {
		appReferences := make([]struct {
			Id   string                                                                  `json:"id"`
			Type openapi.UserInvitationCreateRequestDataRelationshipsVisibleAppsDataType `json:"type"`
		}, len(user.VisibleAppIDs))
		for index, id := range user.VisibleAppIDs {
			appReferences[index].Id = id
			appReferences[index].Type = openapi.UserInvitationCreateRequestDataRelationshipsVisibleAppsDataTypeApps
		}
		requestObject.Data.Relationships = &struct {
			VisibleApps *struct {
				Data *[]struct {
					Id   string                                                                  `json:"id"`
					Type openapi.UserInvitationCreateRequestDataRelationshipsVisibleAppsDataType `json:"type"`
				} `json:"data,omitempty"`
			} `json:"visibleApps,omitempty"`
		}{
			VisibleApps: &struct {
				Data *[]struct {
					Id   string                                                                  `json:"id"`
					Type openapi.UserInvitationCreateRequestDataRelationshipsVisibleAppsDataType `json:"type"`
				} `json:"data,omitempty"`
			}{
				Data: &appReferences,
			},
		}
	}
	if err := json.NewEncoder(body).Encode(requestObject); err != nil {
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

	userResponse := new(openapi.UserInvitationResponse)
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
		username = derefEmail(userResponse.Data.Attributes.Email)
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
