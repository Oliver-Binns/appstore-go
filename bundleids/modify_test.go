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

func TestModifyBundleID_ReturnsErrorForEmptyID(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(http.StatusOK, `{}`)

	result, err := Modify(httpClient, context.Background(), "https://example.com", "", BundleID{Name: "My App"})

	assert.Nil(t, result)
	assert.EqualError(t, err, "bundle ID cannot be empty")
}

func TestModifyBundleID_MakesRequest(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(http.StatusOK, `{"data":{"type":"bundleIds","id":"ABC123"}}`)

	_, _ = Modify(httpClient, context.Background(), "https://example.com", "ABC123", BundleID{Name: "Updated Name"})

	assert.Equal(t, 1, len(httpClient.Requests))
	assert.Equal(t, "PATCH", httpClient.Requests[0].Method)
	assert.Equal(t, "application/json", httpClient.Requests[0].Header.Get("Content-Type"))
	assert.Equal(t, "https://example.com/v1/bundleIds/ABC123", httpClient.Requests[0].URL.String())

	bodyBytes, err := io.ReadAll(httpClient.Requests[0].Body)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"data": {
			"type": "bundleIds",
			"id": "ABC123",
			"attributes": {
				"name": "Updated Name"
			}
		}
	}`, string(bodyBytes))
}

func TestModifyBundleID_ThrowsErrorIfNon200Returned(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(http.StatusNotFound, `{}`)

	result, err := Modify(httpClient, context.Background(), "https://example.com", "ABC123", BundleID{Name: "Updated Name"})

	assert.Nil(t, result)
	assert.EqualError(t, err, "unexpected status code: 404")
}

func TestModifyBundleID_DecodesResponse(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(
		http.StatusOK,
		`{
			"data": {
				"type": "bundleIds",
				"id": "ABC123",
				"attributes": {
					"name": "Updated Name",
					"identifier": "com.example.myapp",
					"platform": "IOS",
					"seedId": "SEED001"
				}
			}
		}`,
	)

	result, err := Modify(
		httpClient, context.Background(), "https://example.com",
		"ABC123", BundleID{Name: "Updated Name"},
	)

	assert.Nil(t, err)
	assert.Equal(t, "ABC123", result.ID)
	assert.Equal(t, "Updated Name", result.Name)
	assert.Equal(t, "com.example.myapp", result.Identifier)
	assert.Equal(t, openapi.IOS, result.Platform)
	assert.Equal(t, "SEED001", result.SeedID)
}
