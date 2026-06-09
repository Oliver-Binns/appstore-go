package devices

import (
	"context"
	"fmt"
	"net/http"

	"github.com/oliver-binns/appstore-go/networking"
	"github.com/oliver-binns/appstore-go/openapi"
)

func FindByUDID(c networking.HTTPClient, ctx context.Context, rawURL string, udid string) (*Device, error) {
	client, err := openapi.NewClientWithResponses(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	filter := []string{udid}
	resp, err := client.DevicesGetCollectionWithResponse(ctx, &openapi.DevicesGetCollectionParams{
		FilterUdid: &filter,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find device by UDID: %w", err)
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
