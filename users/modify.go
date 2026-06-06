package users

import (
	"bytes"
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
	requestObject.Data.Attributes = &openapi.UserUpdateAttributes{
		AllAppsVisible:      ptr.PtrOrNil(user.AllAppsVisible),
		ProvisioningAllowed: ptr.PtrOrNil(user.ProvisioningAllowed),
		Roles:               ptr.SlicePtrOrNil(user.Roles),
	}
	if len(user.VisibleAppIDs) != 0 || !user.AllAppsVisible {
		appReferences := make([]openapi.AppReference, len(user.VisibleAppIDs))
		for index, id := range user.VisibleAppIDs {
			appReferences[index] = openapi.AppReference{
				Id:   id,
				Type: openapi.AppReferenceTypeApps,
			}
		}
		requestObject.Data.Relationships = &openapi.UserUpdateRelationships{
			VisibleApps: &openapi.VisibleAppsRelationship{Data: &appReferences},
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

	return &User{
		ID:                  userResponse.Data.Id,
		FirstName:           ptr.Deref(userResponse.Data.Attributes.FirstName),
		LastName:            ptr.Deref(userResponse.Data.Attributes.LastName),
		Username:            ptr.Deref(userResponse.Data.Attributes.Username),
		Roles:               ptr.Deref(userResponse.Data.Attributes.Roles),
		AllAppsVisible:      ptr.Deref(userResponse.Data.Attributes.AllAppsVisible),
		ProvisioningAllowed: ptr.Deref(userResponse.Data.Attributes.ProvisioningAllowed),
		// Visible App IDs are returned from the input as these are not available in the API response:
		VisibleAppIDs: user.VisibleAppIDs,
	}, nil
}
