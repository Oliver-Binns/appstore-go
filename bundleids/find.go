package bundleids

import (
	"context"
	"fmt"
	"net/http"

	"github.com/oliver-binns/appstore-go/networking"
	"github.com/oliver-binns/appstore-go/openapi"
)

// maxPageSize is the largest page size App Store Connect accepts for collection requests.
const maxPageSize = 200

// FindByIdentifier returns the bundle ID whose identifier exactly matches identifier, or nil if none does.
func FindByIdentifier(c networking.HTTPClient, ctx context.Context, rawURL string, identifier string) (*BundleID, error) {
	client, err := openapi.NewClientWithResponses(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	// The identifier filter also matches longer identifiers (e.g. com.example.app.widget
	// for com.example.app), so request the largest page and pick the exact match.
	filter := []string{identifier}
	limit := maxPageSize
	resp, err := client.BundleIdsGetCollectionWithResponse(ctx, &openapi.BundleIdsGetCollectionParams{
		FilterIdentifier: &filter,
		Limit:            &limit,
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

	for _, data := range resp.JSON200.Data {
		if bundleID := fromResponse(data); bundleID.Identifier == identifier {
			return bundleID, nil
		}
	}
	return nil, nil
}
