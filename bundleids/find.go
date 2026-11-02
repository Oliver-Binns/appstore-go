package bundleids

import (
	"context"
	"fmt"
	"net/http"

	"github.com/oliver-binns/appstore-go/networking"
	"github.com/oliver-binns/appstore-go/openapi"
)

func FindByIdentifier(c networking.HTTPClient, ctx context.Context, rawURL string, identifier string) (*BundleID, error) {
	client, err := openapi.NewClientWithResponses(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	filter := []string{identifier}
	resp, err := client.BundleIdsGetCollectionWithResponse(ctx, &openapi.BundleIdsGetCollectionParams{
		FilterIdentifier: &filter,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find bundle ID by identifier: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("failed to decode response")
	}

	if len(resp.JSON200.Data) == 0 {
		return nil, nil
	}

	return fromResponse(resp.JSON200.Data[0]), nil
}
