package devices

import (
	"context"
	"net/http"
	"testing"

	"github.com/oliver-binns/appstore-go/mocknetworking"
	"github.com/oliver-binns/appstore-go/openapi"
	"github.com/stretchr/testify/assert"
)

func TestFindByUDID_MakesRequest(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWith200Response(`{"data": [], "links": {}}`)

	_, _ = FindByUDID(httpClient, context.Background(), "https://example.com", "00008101-001234AB3C04001E")

	assert.Equal(t, 1, len(httpClient.Requests))
	assert.Equal(t, "GET", httpClient.Requests[0].Method)
	assert.Equal(t, "https://example.com/v1/devices?filter%5Budid%5D=00008101-001234AB3C04001E", httpClient.Requests[0].URL.String())
}

func TestFindByUDID_ReturnsErrorForUnexpectedStatusCode(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(http.StatusUnauthorized, `{}`)

	device, err := FindByUDID(httpClient, context.Background(), "https://example.com", "00008101-001234AB3C04001E")

	assert.Nil(t, device)
	assert.EqualError(t, err, "unexpected status code: 401")
}

func TestFindByUDID_ReturnsNilIfNotFound(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWith200Response(`{"data": [], "links": {}}`)

	device, err := FindByUDID(httpClient, context.Background(), "https://example.com", "00008101-001234AB3C04001E")

	assert.Nil(t, err)
	assert.Nil(t, device)
}

func TestFindByUDID_DecodesResponse(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWith200Response(`{
		"data": [
			{
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
		],
		"links": {}
	}`)

	device, err := FindByUDID(httpClient, context.Background(), "https://example.com", "00008101-001234AB3C04001E")

	assert.Nil(t, err)
	assert.NotNil(t, device)
	assert.Equal(t, "69a495c9-7dbc-5733-e053-5b8c7c1155b0", device.ID)
	assert.Equal(t, "Oliver's iPhone", device.Name)
	assert.Equal(t, "00008101-001234AB3C04001E", device.UDID)
	assert.Equal(t, openapi.IPHONE, device.DeviceClass)
	assert.Equal(t, "iPhone 14 Pro", device.Model)
	assert.Equal(t, openapi.IOS, device.Platform)
	assert.Equal(t, openapi.Enabled, device.Status)
}
