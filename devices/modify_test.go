package devices

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/oliver-binns/appstore-go/mocknetworking"
	"github.com/stretchr/testify/assert"
)

func TestModifyDevice_MakesRequest(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWith200Response(`{}`)

	_, _ = Modify(
		httpClient, context.Background(), "https://example.com", "abcd1234-5678-90ab-cdef-1234567890ab",
		Device{Name: "My iPhone", Status: Enabled},
	)

	assert.Equal(t, 1, len(httpClient.Requests))
	assert.Equal(t, "PATCH", httpClient.Requests[0].Method)
	assert.Equal(t, "application/json", httpClient.Requests[0].Header.Get("Content-Type"))
	assert.Equal(t, "https://example.com/devices/abcd1234-5678-90ab-cdef-1234567890ab", httpClient.Requests[0].URL.String())

	bodyBytes, err := io.ReadAll(httpClient.Requests[0].Body)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"data": {
			"id": "abcd1234-5678-90ab-cdef-1234567890ab",
			"type": "devices",
			"attributes": {
				"name": "My iPhone",
				"status": "ENABLED"
			}
		}
	}`, string(bodyBytes))
}

func TestModifyDevice_ThrowsErrorForUnexpectedStatusCode(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(http.StatusBadRequest, `{}`)

	d, err := Modify(
		httpClient, context.Background(), "https://example.com", "abcd1234-5678-90ab-cdef-1234567890ab",
		Device{Name: "My iPhone", Status: Enabled},
	)

	assert.Nil(t, d)
	assert.EqualError(t, err, "unexpected status code: 400")
}

func TestModifyDevice_DecodesResponse(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWith200Response(`{
		"data": {
			"type": "devices",
			"id": "69a495c9-7dbc-5733-e053-5b8c7c1155b0",
			"attributes": {
				"name": "My iPhone",
				"udid": "00008101-001234AB3C04001E",
				"deviceClass": "IPHONE",
				"model": "iPhone 14 Pro",
				"os": "16.1",
				"platform": "IOS",
				"status": "ENABLED"
			}
		}
	}`)

	device, err := Modify(
		httpClient, context.Background(), "https://example.com", "69a495c9-7dbc-5733-e053-5b8c7c1155b0",
		Device{Name: "My iPhone", Status: Enabled},
	)

	assert.Nil(t, err)
	assert.Equal(t, "69a495c9-7dbc-5733-e053-5b8c7c1155b0", device.ID)
	assert.Equal(t, "My iPhone", device.Name)
	assert.Equal(t, "00008101-001234AB3C04001E", device.UDID)
	assert.Equal(t, IPhone, device.DeviceClass)
	assert.Equal(t, "iPhone 14 Pro", device.Model)
	assert.Equal(t, IOS, device.Platform)
	assert.Equal(t, Enabled, device.Status)
}
