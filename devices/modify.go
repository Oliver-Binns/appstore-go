package devices

import (
	"context"
	"fmt"

	"github.com/oliver-binns/appstore-go/internal/ptr"
	"github.com/oliver-binns/appstore-go/networking"
	"github.com/oliver-binns/appstore-go/openapi"
)

func Modify(c networking.HTTPClient, ctx context.Context, rawURL string, id string, device Device) (*Device, error) {
	if id == "" {
		return nil, fmt.Errorf("device ID cannot be empty")
	}

	client, err := openapi.NewClientWithResponses(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	requestObject := openapi.DeviceUpdateRequest{}
	requestObject.Data.Id = id
	requestObject.Data.Type = openapi.Devices
	requestObject.Data.Attributes = &openapi.DeviceUpdateAttributes{
		Name:   ptr.PtrOrNil(device.Name),
		Status: ptr.PtrOrNil(device.Status),
	}

	resp, err := client.DevicesUpdateInstanceWithResponse(ctx, id, requestObject)
	if err != nil {
		return nil, fmt.Errorf("failed to modify device: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("failed to decode response")
	}

	return fromResponse(resp.JSON200.Data), nil
}
