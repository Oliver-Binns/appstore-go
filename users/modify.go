package users

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/oliver-binns/appstore-go/connectapi"
	"github.com/oliver-binns/googleplay-go/networking"
)

func Modify(c networking.HTTPClient, ctx context.Context, rawURL string, id string, user User) (*User, error) {
	// Parse the raw URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	parsedURL.Path = path.Join(parsedURL.Path, "users", id)

	// Create the request body
	body := bytes.NewBuffer(nil)
	requestObject := connectapi.Request[User]{
		Data: connectapi.RequestData[User]{
			ID:   id,
			Type: "users",
			Data: user,
		},
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

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Response: %+v\n", resp)

		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Printf("Body: %s\n", bodyString)

		return nil, fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}

	if err != nil {
		return nil, err
	}

	userResponse := new(connectapi.Response[User])
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
		Username:            userResponse.Data.Data.Username,
		Roles:               userResponse.Data.Data.Roles,
		AllAppsVisible:      userResponse.Data.Data.AllAppsVisible,
		ProvisioningAllowed: userResponse.Data.Data.ProvisioningAllowed,
	}, nil
}
