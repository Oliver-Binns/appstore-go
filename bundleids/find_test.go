package bundleids

import (
	"context"
	"net/http"
	"testing"

	"github.com/oliver-binns/appstore-go/mocknetworking"
	"github.com/oliver-binns/appstore-go/openapi"
	"github.com/stretchr/testify/assert"
)

func TestFindByIdentifier_MakesRequest(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(http.StatusOK, `{"data":[]}`)

	_, _ = FindByIdentifier(httpClient, context.Background(), "https://example.com", "com.example.myapp")

	assert.Equal(t, 1, len(httpClient.Requests))
	assert.Equal(t, "GET", httpClient.Requests[0].Method)
	assert.Equal(t, "https://example.com/v1/bundleIds?filter%5Bidentifier%5D=com.example.myapp", httpClient.Requests[0].URL.String())
}

func TestFindByIdentifier_ThrowsErrorIfNon200Returned(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(http.StatusBadRequest, `{}`)

	result, err := FindByIdentifier(httpClient, context.Background(), "https://example.com", "com.example.myapp")

	assert.Nil(t, result)
	assert.EqualError(t, err, "unexpected status code: 400")
}

func TestFindByIdentifier_ReturnsNilWhenNotFound(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(http.StatusOK, `{"data":[]}`)

	result, err := FindByIdentifier(httpClient, context.Background(), "https://example.com", "com.example.myapp")

	assert.Nil(t, err)
	assert.Nil(t, result)
}

func TestFindByIdentifier_DecodesResponse(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(
		http.StatusOK,
		`{
			"data": [{
				"type": "bundleIds",
				"id": "ABC123",
				"attributes": {
					"name": "My App",
					"identifier": "com.example.myapp",
					"platform": "IOS",
					"seedId": "SEED001"
				}
			}]
		}`,
	)

	result, err := FindByIdentifier(httpClient, context.Background(), "https://example.com", "com.example.myapp")

	assert.Nil(t, err)
	assert.Equal(t, "ABC123", result.ID)
	assert.Equal(t, "My App", result.Name)
	assert.Equal(t, "com.example.myapp", result.Identifier)
	assert.Equal(t, openapi.IOS, result.Platform)
	assert.Equal(t, "SEED001", result.SeedID)
}
