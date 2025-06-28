package users

import (
	"context"
	"net/http"
	"testing"

	"github.com/oliver-binns/appstore-go/mocknetworking"
	"github.com/stretchr/testify/assert"
)

func TestGetUser_MakesRequest(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWith200Response(`{ }`)

	_, _ = Get(
		httpClient, context.Background(), "https://example.com", "abcd1234-5678-90ab-cdef-1234567890ab",
	)

	assert.Equal(t, len(httpClient.Requests), 1)
	assert.Equal(t, httpClient.Requests[0].Method, "GET")
	assert.Equal(t, httpClient.Requests[0].URL.String(), "https://example.com/users/abcd1234-5678-90ab-cdef-1234567890ab")
}

func TestGetUser_MakesSecondRequestToInvitations_When404Returned(t *testing.T) {
	notFound := http.StatusNotFound

	httpClient := mocknetworking.MockHTTPClient{
		Responses: []mocknetworking.MockHTTPResponse{
			{
				StatusCode: &notFound,
				Body:       `{ }`,
			},
			{
				Body: `{ }`,
			},
		},
	}

	_, _ = Get(
		&httpClient, context.Background(), "https://example.com", "abcd1234-5678-90ab-cdef-1234567890ab",
	)

	assert.Equal(t, len(httpClient.Requests), 2)
	assert.Equal(t, httpClient.Requests[0].Method, "GET")
	assert.Equal(t, httpClient.Requests[0].URL.String(), "https://example.com/users/abcd1234-5678-90ab-cdef-1234567890ab")

	assert.Equal(t, httpClient.Requests[1].Method, "GET")
	assert.Equal(t, httpClient.Requests[1].URL.String(), "https://example.com/userInvitations/abcd1234-5678-90ab-cdef-1234567890ab")
}

func TestGetUser_DecodesResponse(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWith200Response(`{
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
	)

	user, err := Get(
		httpClient, context.Background(), "https://example.com", "abc",
	)

	assert.Nil(t, err)

	assert.Equal(t, user.ID, "69a495c9-7dbc-5733-e053-5b8c7c1155b0")
	assert.Equal(t, user.FirstName, "Oliver")
	assert.Equal(t, user.LastName, "Binns")
	assert.Equal(t, user.Username, "mail@oliverbinns.co.uk")
}

func TestGetUser_DecodesInvitationResponse_When404ReturnedFromUsers(t *testing.T) {
	notFound := http.StatusNotFound

	httpClient := mocknetworking.MockHTTPClient{
		Responses: []mocknetworking.MockHTTPResponse{
			{
				StatusCode: &notFound,
				Body:       `{ }`,
			},
			{
				Body: `
				{
					"data": {
						"type": "users",
						"id": "69a495c9-7dbc-5733-e053-5b8c7c1155b0",
						"attributes": {
							"allAppsVisible": true,
							"lastName": "Binns",
							"firstName": "Oliver",
							"provisioningAllowed": true,
							"roles": ["ACCOUNT_HOLDER", "ADMIN"],
							"email": "mail@oliverbinns.co.uk"
						}
					}
				}`,
			},
		},
	}

	user, err := Get(
		&httpClient, context.Background(), "https://example.com", "abcd1234-5678-90ab-cdef-1234567890ab",
	)

	assert.Nil(t, err)

	assert.Equal(t, user.ID, "69a495c9-7dbc-5733-e053-5b8c7c1155b0")
	assert.Equal(t, user.FirstName, "Oliver")
	assert.Equal(t, user.LastName, "Binns")
	assert.Equal(t, user.Username, "mail@oliverbinns.co.uk")
}
