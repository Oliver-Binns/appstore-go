package bundleids

import (
	"context"
	"net/http"
	"testing"

	"github.com/oliver-binns/appstore-go/mocknetworking"
	"github.com/stretchr/testify/assert"
)

func TestDeleteBundleID_MakesRequest(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(http.StatusNoContent, ``)

	_ = Delete(httpClient, context.Background(), "https://example.com", "ABC123")

	assert.Equal(t, 1, len(httpClient.Requests))
	assert.Equal(t, "DELETE", httpClient.Requests[0].Method)
	assert.Equal(t, "https://example.com/v1/bundleIds/ABC123", httpClient.Requests[0].URL.String())
}

func TestDeleteBundleID_ReturnsNilOnSuccess(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(http.StatusNoContent, ``)

	err := Delete(httpClient, context.Background(), "https://example.com", "ABC123")

	assert.Nil(t, err)
}

func TestDeleteBundleID_ThrowsErrorIfNon204Returned(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(http.StatusNotFound, `{}`)

	err := Delete(httpClient, context.Background(), "https://example.com", "ABC123")

	assert.EqualError(t, err, "unexpected response code: 404")
}
