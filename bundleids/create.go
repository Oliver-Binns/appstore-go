package bundleids

import (
	"context"
	"fmt"

	"github.com/oliver-binns/appstore-go/internal/ptr"
	"github.com/oliver-binns/appstore-go/networking"
	"github.com/oliver-binns/appstore-go/openapi"
)

func Create(c networking.HTTPClient, ctx context.Context, rawURL string, bundleID BundleID) (*BundleID, error) {
	client, err := openapi.NewClientWithResponses(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	requestObject := openapi.BundleIdCreateRequest{}
	requestObject.Data.Type = openapi.BundleIdCreateRequestDataTypeBundleIds
	requestObject.Data.Attributes.Identifier = bundleID.Identifier
	requestObject.Data.Attributes.Name = bundleID.Name
	requestObject.Data.Attributes.Platform = bundleID.Platform
	requestObject.Data.Attributes.SeedId = ptr.PtrOrNil(bundleID.SeedID)

	resp, err := client.BundleIdsCreateInstanceWithResponse(ctx, requestObject)
	if err != nil {
		return nil, fmt.Errorf("failed to create bundle ID: %w", err)
	}
	if resp.StatusCode() != 201 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	if resp.JSON201 == nil {
		return nil, fmt.Errorf("failed to decode response")
	}

	return fromResponse(resp.JSON201.Data), nil
}
