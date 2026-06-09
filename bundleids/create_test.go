package bundleids

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/oliver-binns/appstore-go/mocknetworking"
	"github.com/oliver-binns/appstore-go/openapi"
	"github.com/stretchr/testify/assert"
)

func TestCreateBundleID_MakesRequest(t *testing.T) {
	bundleID := BundleID{
		Name:       "My App",
		Identifier: "com.example.myapp",
		Platform:   openapi.IOS,
	}
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(http.StatusCreated, `{"data":{"type":"bundleIds","id":"ABC123"}}`)

	_, _ = Create(httpClient, context.Background(), "https://example.com", bundleID)

	assert.Equal(t, 1, len(httpClient.Requests))
	assert.Equal(t, "POST", httpClient.Requests[0].Method)
	assert.Equal(t, "application/json", httpClient.Requests[0].Header.Get("Content-Type"))
	assert.Equal(t, "https://example.com/v1/bundleIds", httpClient.Requests[0].URL.String())

	bodyBytes, err := io.ReadAll(httpClient.Requests[0].Body)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"data": {
			"type": "bundleIds",
			"attributes": {
				"identifier": "com.example.myapp",
				"name": "My App",
				"platform": "IOS"
			}
		}
	}`, string(bodyBytes))
}

func TestCreateBundleID_ThrowsErrorIfNon201Returned(t *testing.T) {
	bundleID := BundleID{
		Name:       "My App",
		Identifier: "com.example.myapp",
		Platform:   openapi.IOS,
	}
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(http.StatusBadRequest, `{}`)

	result, err := Create(httpClient, context.Background(), "https://example.com", bundleID)

	assert.Nil(t, result)
	assert.EqualError(t, err, "unexpected status code: 400")
}

func TestCreateBundleID_DecodesResponse(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(
		http.StatusCreated,
		`{
			"data": {
				"type": "bundleIds",
				"id": "ABC123",
				"attributes": {
					"name": "My App",
					"identifier": "com.example.myapp",
					"platform": "IOS",
					"seedId": "SEED001"
				}
			}
		}`,
	)

	result, err := Create(
		httpClient, context.Background(), "https://example.com",
		BundleID{
			Name:       "My App",
			Identifier: "com.example.myapp",
			Platform:   openapi.IOS,
		},
	)

	assert.Nil(t, err)
	assert.Equal(t, "ABC123", result.ID)
	assert.Equal(t, "My App", result.Name)
	assert.Equal(t, "com.example.myapp", result.Identifier)
	assert.Equal(t, openapi.IOS, result.Platform)
	assert.Equal(t, "SEED001", result.SeedID)
}
