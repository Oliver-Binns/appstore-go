package users

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUser_MakesRequest(t *testing.T) {
	httpClient := &mockHTTPClient{
		response: `{
			"users": [
			
			],
			"nextPageToken": string
		}`,
	}

	_, _ = Get(
		httpClient, context.Background(), "https://example.com", "abcd1234-5678-90ab-cdef-1234567890ab",
	)

	assert.Equal(t, len(httpClient.requests), 1)
	assert.Equal(t, httpClient.requests[0].Method, "GET")
	assert.Equal(t, httpClient.requests[0].URL.String(), "https://example.com/users/abcd1234-5678-90ab-cdef-1234567890ab")
}

func TestGetUser_DecodesResponse(t *testing.T) {
	httpClient := &mockHTTPClient{
		response: `{
  "data" : {
    "type" : "users",
    "id" : "69a495c9-7dbc-5733-e053-5b8c7c1155b0",
    "attributes" : {
      "allAppsVisible" : true,
      "lastName" : "Binns",
      "firstName" : "Oliver",
      "provisioningAllowed" : true,
      "roles" : [ "ACCOUNT_HOLDER", "ADMIN" ],
      "username" : "mail@oliverbinns.co.uk"
    },
    "relationships" : {
      "visibleApps" : {
        "links" : {
          "self" : "https://api.appstoreconnect.apple.com/v1/users/69a495c9-7dbc-5733-e053-5b8c7c1155b0/relationships/visibleApps",
          "related" : "https://api.appstoreconnect.apple.com/v1/users/69a495c9-7dbc-5733-e053-5b8c7c1155b0/visibleApps"
        }
      }
    },
    "links" : {
      "self" : "https://api.appstoreconnect.apple.com/v1/users/69a495c9-7dbc-5733-e053-5b8c7c1155b0"
    }
  },
  "links" : {
    "self" : "https://api.appstoreconnect.apple.com/v1/users/69a495c9-7dbc-5733-e053-5b8c7c1155b0"
  }
}`,
	}

	user, _ := Get(
		httpClient, context.Background(), "https://example.com", "abc",
	)

	assert.Equal(t, user.ID, "69a495c9-7dbc-5733-e053-5b8c7c1155b0")
	assert.Equal(t, user.FirstName, "Oliver")
	assert.Equal(t, user.LastName, "Binns")
	assert.Equal(t, user.Username, "mail@oliverbinns.co.uk")
}

type mockHTTPClient struct {
	requests []*http.Request
	response string

	statusCode *int
}

func (c *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.requests = append(c.requests, req)

	responseBody := io.NopCloser(bytes.NewReader([]byte(c.response)))

	status := http.StatusOK
	if c.statusCode != nil {
		status = *c.statusCode
	}

	return &http.Response{
		StatusCode: status,
		Body:       responseBody,
	}, nil
}
