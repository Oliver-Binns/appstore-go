package devices

import (
	"context"
	"net/http"
	"testing"

	"github.com/oliver-binns/appstore-go/mocknetworking"
	"github.com/oliver-binns/appstore-go/openapi"
	"github.com/stretchr/testify/assert"
)

func TestGetDevice_ReturnsErrorForEmptyID(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWith200Response(`{}`)

	_, err := Get(
		httpClient, context.Background(), "https://example.com", "",
	)

	assert.ErrorContains(t, err, "device ID cannot be empty")
	assert.Equal(t, 0, len(httpClient.Requests))
}

func TestGetDevice_MakesRequest(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWith200Response(`{"data":{"type":"devices","id":"abcd1234-5678-90ab-cdef-1234567890ab"}}`)

	_, _ = Get(
		httpClient, context.Background(), "https://example.com", "abcd1234-5678-90ab-cdef-1234567890ab",
	)

	assert.Equal(t, 1, len(httpClient.Requests))
	assert.Equal(t, "GET", httpClient.Requests[0].Method)
	assert.Equal(t, "https://example.com/v1/devices/abcd1234-5678-90ab-cdef-1234567890ab", httpClient.Requests[0].URL.String())
}

func TestGetDevice_ThrowsErrorForUnexpectedStatusCode(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(http.StatusNotFound, `{}`)

	d, err := Get(
		httpClient, context.Background(), "https://example.com", "abcd1234-5678-90ab-cdef-1234567890ab",
	)

	assert.Nil(t, d)
	assert.EqualError(t, err, "unexpected status code: 404")
}

func TestGetDevice_DecodesResponse(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWith200Response(`{
		"data": {
			"type": "devices",
			"id": "69a495c9-7dbc-5733-e053-5b8c7c1155b0",
			"attributes": {
				"name": "Oliver's iPhone",
				"udid": "00008101-001234AB3C04001E",
				"deviceClass": "IPHONE",
				"model": "iPhone 14 Pro",
				"platform": "IOS",
				"status": "ENABLED"
			}
		}
	}`)

	device, err := Get(
		httpClient, context.Background(), "https://example.com", "69a495c9-7dbc-5733-e053-5b8c7c1155b0",
	)

	assert.Nil(t, err)
	assert.Equal(t, "69a495c9-7dbc-5733-e053-5b8c7c1155b0", device.ID)
	assert.Equal(t, "Oliver's iPhone", device.Name)
	assert.Equal(t, "00008101-001234AB3C04001E", device.UDID)
	assert.Equal(t, openapi.IPHONE, device.DeviceClass)
	assert.Equal(t, "iPhone 14 Pro", device.Model)
	assert.Equal(t, openapi.IOS, device.Platform)
	assert.Equal(t, openapi.Enabled, device.Status)
}
