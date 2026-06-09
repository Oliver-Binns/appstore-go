package bundleids

import (
	"context"
	"fmt"

	"github.com/oliver-binns/appstore-go/networking"
	"github.com/oliver-binns/appstore-go/openapi"
)

func Get(c networking.HTTPClient, ctx context.Context, rawURL string, id string) (*BundleID, error) {
	if id == "" {
		return nil, fmt.Errorf("bundle ID cannot be empty")
	}

	client, err := openapi.NewClientWithResponses(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	resp, err := client.BundleIdsGetInstanceWithResponse(ctx, id, &openapi.BundleIdsGetInstanceParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to get bundle ID: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("failed to decode response")
	}

	return fromResponse(resp.JSON200.Data), nil
}
