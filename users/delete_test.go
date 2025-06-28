package users

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/oliver-binns/appstore-go/mocknetworking"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUser_MakesRequest(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWith200Response(`{ }`)

	_ = Delete(
		httpClient, context.Background(), "https://example.com", "user-id",
	)

	assert.Equal(t, len(httpClient.Requests), 1)
	assert.Equal(t, httpClient.Requests[0].Method, "DELETE")
	assert.Equal(t, httpClient.Requests[0].URL.String(), "https://example.com/users/user-id")

	// The body should be empty for a DELETE request
	bodyBytes, err := io.ReadAll(httpClient.Requests[0].Body)
	assert.NoError(t, err)
	assert.Equal(t, len(bodyBytes), 0)
}

func TestDeleteUser_ReturnsNilForSuccess(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(http.StatusNoContent, `{ }`)

	err := Delete(
		httpClient, context.Background(), "https://example.com", "user-id",
	)

	assert.NoError(t, err)
}

func TestDeleteUser_ThrowsErrorForUnexpectedStatusCode(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWith200Response(`{ }`)

	err := Delete(
		httpClient, context.Background(), "https://example.com", "user-id",
	)

	assert.Equal(t, err.Error(), "unexpected response code: 200")
}
