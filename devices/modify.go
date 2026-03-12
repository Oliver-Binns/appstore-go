package devices

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

func Modify(c networking.HTTPClient, ctx context.Context, rawURL string, id string, device Device) (*Device, error) {
	// Parse the raw URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	parsedURL.Path = path.Join(parsedURL.Path, "devices", id)

	requestData := openapi.DeviceUpdateRequestData{
		Id:   id,
		Type: "devices",
		Attributes: openapi.DeviceUpdateRequestAttributes{
			Name:   stringPtrOrNil(device.Name),
			Status: statusPtrOrNil(device.Status),
		},
	}

	// Create the request body
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(openapi.DeviceUpdateRequest{Data: requestData}); err != nil {
		return nil, fmt.Errorf("failed to encode request body: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, parsedURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var deviceResponse openapi.DeviceResponse
	if err := json.NewDecoder(resp.Body).Decode(&deviceResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if err := resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close response body: %w", err)
	}

	return fromResponse(deviceResponse.Data), nil
}

// stringPtrOrNil returns a pointer to s when s is non-empty and nil otherwise,
// so that empty strings are omitted from the JSON request body.
func stringPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// statusPtrOrNil returns a pointer to s when s is non-empty and nil otherwise.
func statusPtrOrNil(s openapi.DeviceStatus) *openapi.DeviceStatus {
	if s == "" {
		return nil
	}
	return &s
}
