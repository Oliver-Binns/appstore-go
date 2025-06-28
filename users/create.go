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

func Create(c networking.HTTPClient, ctx context.Context, rawURL string, user UserInvitation) (*UserInvitation, error) {
	// Parse the raw URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	parsedURL.Path = path.Join(parsedURL.Path, "userInvitations")

	// Create the request body
	body := bytes.NewBuffer(nil)
	requestObject := connectapi.Request[UserInvitation]{
		Data: connectapi.RequestData[UserInvitation]{
			Type: "userInvitations",
			Data: user,
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
	}

	userResponse := new(connectapi.Response[UserInvitation])
	if err := json.NewDecoder(resp.Body).Decode(userResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if err := resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close response body: %w", err)
	}

	userResponse.Data.Data.ID = userResponse.Data.ID
	return &userResponse.Data.Data, nil
}
