package devices

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/oliver-binns/appstore-go/mocknetworking"
	"github.com/stretchr/testify/assert"
)

func TestRegisterDevice_MakesRequest(t *testing.T) {
	device := Device{
		Name:     "Oliver's iPhone",
		UDID:     "00008101-001234AB3C04001E",
		Platform: IOS,
	}
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(http.StatusCreated, `{}`)

	_, _ = Register(
		httpClient, context.Background(), "https://example.com", device,
	)

	assert.Equal(t, 1, len(httpClient.Requests))
	assert.Equal(t, "POST", httpClient.Requests[0].Method)
	assert.Equal(t, "application/json", httpClient.Requests[0].Header.Get("Content-Type"))
	assert.Equal(t, "https://example.com/devices", httpClient.Requests[0].URL.String())

	bodyBytes, err := io.ReadAll(httpClient.Requests[0].Body)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"data": {
			"type": "devices",
			"attributes": {
				"name": "Oliver's iPhone",
				"udid": "00008101-001234AB3C04001E",
				"platform": "IOS"
			}
		}
	}`, string(bodyBytes))
}

func TestRegisterDevice_ThrowsErrorIfNon201Returned(t *testing.T) {
	device := Device{
		Name:     "Oliver's iPhone",
		UDID:     "00008101-001234AB3C04001E",
		Platform: IOS,
	}
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(http.StatusBadRequest, `{}`)

	d, err := Register(
		httpClient, context.Background(), "https://example.com", device,
	)

	assert.Nil(t, d)
	assert.EqualError(t, err, "unexpected status code: 400")
}

func TestRegisterDevice_DecodesResponse(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(
		http.StatusCreated,
		`{
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
		}`,
	)

	device, err := Register(
		httpClient, context.Background(), "https://example.com",
		Device{
			Name:     "Oliver's iPhone",
			UDID:     "00008101-001234AB3C04001E",
			Platform: IOS,
		},
	)

	assert.Nil(t, err)
	assert.Equal(t, "69a495c9-7dbc-5733-e053-5b8c7c1155b0", device.ID)
	assert.Equal(t, "Oliver's iPhone", device.Name)
	assert.Equal(t, "00008101-001234AB3C04001E", device.UDID)
	assert.Equal(t, IPhone, device.DeviceClass)
	assert.Equal(t, "iPhone 14 Pro", device.Model)
	assert.Equal(t, IOS, device.Platform)
	assert.Equal(t, Enabled, device.Status)
}
