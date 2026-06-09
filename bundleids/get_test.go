package bundleids

import (
	"context"
	"net/http"
	"testing"

	"github.com/oliver-binns/appstore-go/mocknetworking"
	"github.com/oliver-binns/appstore-go/openapi"
	"github.com/stretchr/testify/assert"
)

func TestGetBundleID_ReturnsErrorForEmptyID(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(http.StatusOK, `{}`)

	result, err := Get(httpClient, context.Background(), "https://example.com", "")

	assert.Nil(t, result)
	assert.EqualError(t, err, "bundle ID cannot be empty")
}

func TestGetBundleID_MakesRequest(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(http.StatusOK, `{"data":{"type":"bundleIds","id":"ABC123"}}`)

	_, _ = Get(httpClient, context.Background(), "https://example.com", "ABC123")

	assert.Equal(t, 1, len(httpClient.Requests))
	assert.Equal(t, "GET", httpClient.Requests[0].Method)
	assert.Equal(t, "https://example.com/v1/bundleIds/ABC123", httpClient.Requests[0].URL.String())
}

func TestGetBundleID_ThrowsErrorIfNon200Returned(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(http.StatusNotFound, `{}`)

	result, err := Get(httpClient, context.Background(), "https://example.com", "ABC123")

	assert.Nil(t, result)
	assert.EqualError(t, err, "unexpected status code: 404")
}

func TestGetBundleID_DecodesResponse(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(
		http.StatusOK,
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

	result, err := Get(httpClient, context.Background(), "https://example.com", "ABC123")

	assert.Nil(t, err)
	assert.Equal(t, "ABC123", result.ID)
	assert.Equal(t, "My App", result.Name)
	assert.Equal(t, "com.example.myapp", result.Identifier)
	assert.Equal(t, openapi.IOS, result.Platform)
	assert.Equal(t, "SEED001", result.SeedID)
}
