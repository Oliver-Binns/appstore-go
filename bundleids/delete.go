package bundleids

import (
	"context"
	"fmt"
	"net/http"

	"github.com/oliver-binns/appstore-go/networking"
	"github.com/oliver-binns/appstore-go/openapi"
)

func Delete(c networking.HTTPClient, ctx context.Context, rawURL string, id string) error {
	client, err := openapi.NewClientWithResponses(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	resp, err := client.BundleIdsDeleteInstanceWithResponse(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete bundle ID: %w", err)
	}

	if resp.StatusCode() == http.StatusNoContent {
		return nil
	}
	return fmt.Errorf("unexpected response code: %d", resp.StatusCode())
}
