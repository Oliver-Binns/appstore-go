package users

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser_MakesRequest(t *testing.T) {
	user := User{
		FirstName: "Joseph",
		LastName:  "Bloggs",
		Username:  "joe.bloggs@example.com",
		Roles:     []UserRole{Marketing},
	}
	httpClient := &mockHTTPClient{
		response: ``,
	}

	_, _ = Create(
		httpClient, context.Background(), "https://example.com", user,
	)

	assert.Equal(t, len(httpClient.requests), 1)
	assert.Equal(t, httpClient.requests[0].Method, "POST")
	assert.Equal(t, httpClient.requests[0].Header.Get("Content-Type"), "application/json")
	assert.Equal(t, httpClient.requests[0].URL.String(), "https://example.com/userInvitations")

	bodyBytes, err := io.ReadAll(httpClient.requests[0].Body)
	assert.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `{"data":{"type":"userInvitations","attributes":{"firstName":"Joseph","lastName":"Bloggs","username":"joe.bloggs@example.com","roles":["MARKETING"]}}}
`, bodyString)
}

func TestCreateUser_DecodesResponse(t *testing.T) {
	httpClient := &mockHTTPClient{
		response: `{
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
				},
				"relationships": {
					"visibleApps": {
						"links": {
							"self": "https://api.appstoreconnect.apple.com/v1/users/69a495c9-7dbc-5733-e053-5b8c7c1155b0/relationships/visibleApps",
							"related": "https://api.appstoreconnect.apple.com/v1/users/69a495c9-7dbc-5733-e053-5b8c7c1155b0/visibleApps"
						}
					}
				},
				"links": {
					"self": "https://api.appstoreconnect.apple.com/v1/users/69a495c9-7dbc-5733-e053-5b8c7c1155b0"
				}
			},
			"links": {
				"self": "https://api.appstoreconnect.apple.com/v1/users/69a495c9-7dbc-5733-e053-5b8c7c1155b0"
			}
		}`,
	}

	user, _ := Create(
		httpClient, context.Background(), "https://example.com", User{},
	)

	assert.Equal(t, user.FirstName, "Oliver")
	assert.Equal(t, user.LastName, "Binns")
	assert.Equal(t, user.Username, "mail@oliverbinns.co.uk")
}
