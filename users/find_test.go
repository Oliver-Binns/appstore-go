package users

import (
	"context"
	"net/http"
	"testing"

	"github.com/oliver-binns/appstore-go/mocknetworking"
	"github.com/stretchr/testify/assert"
)

func TestFindByEmail_MakesRequest(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWith200Response(`{"data": [], "links": {}}`)

	_, _ = FindByEmail(httpClient, context.Background(), "https://example.com", "mail@oliverbinns.co.uk")

	assert.Equal(t, 1, len(httpClient.Requests))
	assert.Equal(t, "GET", httpClient.Requests[0].Method)
	assert.Equal(t, "https://example.com/users?filter%5Busername%5D=mail%40oliverbinns.co.uk&include=visibleApps", httpClient.Requests[0].URL.String())
}

func TestFindByEmail_ReturnsErrorIfNon200Returned(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(http.StatusUnauthorized, `{}`)

	user, err := FindByEmail(httpClient, context.Background(), "https://example.com", "mail@oliverbinns.co.uk")

	assert.Nil(t, user)
	assert.Equal(t, "unexpected status code: 401", err.Error())
}

func TestFindByEmail_ReturnsNilIfNotFound(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWith200Response(`{"data": [], "links": {}}`)

	user, err := FindByEmail(httpClient, context.Background(), "https://example.com", "notfound@example.com")

	assert.Nil(t, err)
	assert.Nil(t, user)
}

func TestFindByEmail_DecodesResponse(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWith200Response(`{
		"data": [
			{
				"type": "users",
				"id": "69a495c9-7dbc-5733-e053-5b8c7c1155b0",
				"attributes": {
					"allAppsVisible": true,
					"lastName": "Binns",
					"firstName": "Oliver",
					"provisioningAllowed": true,
					"roles": ["ACCOUNT_HOLDER", "ADMIN"],
					"username": "mail@oliverbinns.co.uk"
				},
				"relationships": {
					"visibleApps": {
						"data": [{"id": "123456", "type": "apps"}]
					}
				}
			}
		],
		"links": {}
	}`)

	user, err := FindByEmail(httpClient, context.Background(), "https://example.com", "mail@oliverbinns.co.uk")

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "69a495c9-7dbc-5733-e053-5b8c7c1155b0", user.ID)
	assert.Equal(t, "Oliver", user.FirstName)
	assert.Equal(t, "Binns", user.LastName)
	assert.Equal(t, "mail@oliverbinns.co.uk", user.Username)
	assert.True(t, user.HasAcceptedInvite)
	assert.Equal(t, []string{"123456"}, user.VisibleAppIDs)
}
