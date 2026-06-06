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
	"github.com/oliver-binns/appstore-go/internal/ptr"
	"github.com/oliver-binns/appstore-go/openapi"
	"github.com/oliver-binns/appstore-go/networking"
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
	requestObject.Data.Attributes = openapi.UserInvitationAttributes{
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		Email:               openapi_types.Email(user.Username),
		Roles:               user.Roles,
		AllAppsVisible:      ptr.PtrOrNil(user.AllAppsVisible),
		ProvisioningAllowed: ptr.PtrOrNil(user.ProvisioningAllowed),
	}
	requestObject.Data.Relationships = invitationCreateRelationships(user.VisibleAppIDs, user.AllAppsVisible)
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

	return &User{
		ID:                  userResponse.Data.Id,
		FirstName:           ptr.Deref(userResponse.Data.Attributes.FirstName),
		LastName:            ptr.Deref(userResponse.Data.Attributes.LastName),
		Username:            string(ptr.Deref(userResponse.Data.Attributes.Email)),
		Roles:               ptr.Deref(userResponse.Data.Attributes.Roles),
		AllAppsVisible:      ptr.Deref(userResponse.Data.Attributes.AllAppsVisible),
		ProvisioningAllowed: ptr.Deref(userResponse.Data.Attributes.ProvisioningAllowed),
		// Visible App IDs are returned from the input as these are not available in the API response:
		VisibleAppIDs: user.VisibleAppIDs,
	}, nil
}
