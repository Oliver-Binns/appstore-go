package devices

import (
	"context"
	"fmt"

	"github.com/oliver-binns/appstore-go/networking"
	"github.com/oliver-binns/appstore-go/openapi"
)

func Get(c networking.HTTPClient, ctx context.Context, rawURL string, id string) (*Device, error) {
	if id == "" {
		return nil, fmt.Errorf("device ID cannot be empty")
	}

	client, err := openapi.NewClientWithResponses(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	resp, err := client.DevicesGetInstanceWithResponse(ctx, id, &openapi.DevicesGetInstanceParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("failed to decode response")
	}

	return fromResponse(resp.JSON200.Data), nil
}
