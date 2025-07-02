package users

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/oliver-binns/appstore-go/connectapi"
	"github.com/oliver-binns/googleplay-go/networking"
)

func Create(c networking.HTTPClient, ctx context.Context, rawURL string, user User) (*User, error) {
	// Parse the raw URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	parsedURL.Path = path.Join(parsedURL.Path, "userInvitations")

	invitation := userInvitation{
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		Email:               user.Username,
		Roles:               user.Roles,
		AllAppsVisible:      user.AllAppsVisible,
		ProvisioningAllowed: user.ProvisioningAllowed,
	}

	// Create the request body
	body := bytes.NewBuffer(nil)
	requestObject := connectapi.Request[userInvitation, *userRelationships]{
		Data: connectapi.RequestData[userInvitation, *userRelationships]{
			Type:          "userInvitations",
			Data:          invitation,
			Relationships: user.relationships(),
		},
	}
	err = json.NewEncoder(body).Encode(requestObject)
	if err != nil {
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

	userResponse := new(connectapi.Response[userInvitation, *userRelationships])
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
		// Visible App IDs are returned from the input as these are not available in the API response:
		VisibleAppIDs: user.relationships().ids(),
	}, nil
}
