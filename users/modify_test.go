package users

import (
	"context"
	"io"
	"testing"

	"github.com/oliver-binns/appstore-go/mocknetworking"
	"github.com/stretchr/testify/assert"
)

func TestModifyUser_MakesRequest(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWith200Response(`{ }`)

	_, _ = Modify(
		httpClient, context.Background(), "https://example.com", "abcd1234-5678-90ab-cdef-1234567890ab",
		User{
			AllAppsVisible:      true,
			ProvisioningAllowed: true,
		},
	)

	assert.Equal(t, len(httpClient.Requests), 1)
	assert.Equal(t, httpClient.Requests[0].Method, "PATCH")
	assert.Equal(t, httpClient.Requests[0].Header.Get("Content-Type"), "application/json")
	assert.Equal(t, httpClient.Requests[0].URL.String(), "https://example.com/users/abcd1234-5678-90ab-cdef-1234567890ab")

	bodyBytes, err := io.ReadAll(httpClient.Requests[0].Body)
	assert.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `{"data":{"id":"abcd1234-5678-90ab-cdef-1234567890ab","type":"users","attributes":{"allAppsVisible":true,"provisioningAllowed":true}}}
`, bodyString)
}

func TestModifyUser_DecodesResponse(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWith200Response(`
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
				"username": "mail@oliverbinns.co.uk"
			}
		}
	}`)

	user, _ := Modify(
		httpClient, context.Background(), "https://example.com", "dummy-id",
		User{
			AllAppsVisible:      true,
			ProvisioningAllowed: true,
		},
	)

	assert.Equal(t, user.FirstName, "Oliver")
	assert.Equal(t, user.LastName, "Binns")
	assert.Equal(t, user.Username, "mail@oliverbinns.co.uk")
}
