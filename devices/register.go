package devices

import (
	"context"
	"fmt"

	"github.com/oliver-binns/appstore-go/networking"
	"github.com/oliver-binns/appstore-go/openapi"
)

func Register(c networking.HTTPClient, ctx context.Context, rawURL string, device Device) (*Device, error) {
	client, err := openapi.NewClientWithResponses(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	requestObject := openapi.DeviceCreateRequest{}
	requestObject.Data.Type = openapi.DeviceCreateRequestDataTypeDevices
	requestObject.Data.Attributes.Name = device.Name
	requestObject.Data.Attributes.Udid = device.UDID
	requestObject.Data.Attributes.Platform = device.Platform

	resp, err := client.DevicesCreateInstanceWithResponse(ctx, requestObject)
	if err != nil {
		return nil, fmt.Errorf("failed to register device: %w", err)
	}
	if resp.StatusCode() != 201 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	if resp.JSON201 == nil {
		return nil, fmt.Errorf("failed to decode response")
	}

	return fromResponse(resp.JSON201.Data), nil
}
