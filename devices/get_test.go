package devices

import (
	"context"
	"net/http"
	"testing"

	"github.com/oliver-binns/appstore-go/mocknetworking"
	"github.com/stretchr/testify/assert"
)

func TestGetDevice_MakesRequest(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWith200Response(`{}`)

	_, _ = Get(
		httpClient, context.Background(), "https://example.com", "abcd1234-5678-90ab-cdef-1234567890ab",
	)

	assert.Equal(t, 1, len(httpClient.Requests))
	assert.Equal(t, "GET", httpClient.Requests[0].Method)
	assert.Equal(t, "https://example.com/devices/abcd1234-5678-90ab-cdef-1234567890ab", httpClient.Requests[0].URL.String())
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
				"os": "16.1",
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
	assert.Equal(t, IPhone, device.DeviceClass)
	assert.Equal(t, "iPhone 14 Pro", device.Model)
	assert.Equal(t, "16.1", device.OS)
	assert.Equal(t, IOS, device.Platform)
	assert.Equal(t, Enabled, device.Status)
}
