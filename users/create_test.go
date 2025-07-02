package users

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/oliver-binns/appstore-go/mocknetworking"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_MakesRequest(t *testing.T) {
	user := User{
		FirstName:     "Joseph",
		LastName:      "Bloggs",
		Username:      "joe.bloggs@example.com",
		Roles:         []UserRole{Marketing},
		VisibleAppIDs: []string{"abcd"},
	}
	httpClient := mocknetworking.MockHTTPClientWith200Response(`{ }`)

	_, _ = Create(
		httpClient, context.Background(), "https://example.com", user,
	)

	assert.Equal(t, len(httpClient.Requests), 1)
	assert.Equal(t, httpClient.Requests[0].Method, "POST")
	assert.Equal(t, httpClient.Requests[0].Header.Get("Content-Type"), "application/json")
	assert.Equal(t, httpClient.Requests[0].URL.String(), "https://example.com/userInvitations")

	bodyBytes, err := io.ReadAll(httpClient.Requests[0].Body)
	assert.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `{"data":{"type":"userInvitations","attributes":{"firstName":"Joseph","lastName":"Bloggs","email":"joe.bloggs@example.com","roles":["MARKETING"]},"relationships":{"visibleApps":{"data":[{"id":"abcd","type":"apps"}]}}}}
`, bodyString)
}

func TestCreateUser_ThrowsErrorIfNon201Returned(t *testing.T) {
	user := User{
		FirstName:     "Joseph",
		LastName:      "Bloggs",
		Username:      "joe.bloggs@example.com",
		Roles:         []UserRole{Marketing},
		VisibleAppIDs: []string{"abcd"},
	}
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(
		http.StatusBadRequest,
		`{ }`,
	)

	u, err := Create(
		httpClient, context.Background(), "https://example.com", user,
	)

	assert.Nil(t, u)
	assert.Equal(t, "unexpected status code: 400", err.Error())
}

func TestCreateUser_DecodesResponse(t *testing.T) {
	httpClient := mocknetworking.MockHTTPClientWithSingleResponse(
		http.StatusCreated,
		`{
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
			},
			"relationships": {
			    "visibleApps": {
				    "data": [{
				        "id": "abcdef",
					    "type": "apps"
				    }]
				}
			}
		}
	}`)

	user, err := Create(
		httpClient, context.Background(), "https://example.com", User{
			VisibleAppIDs: []string{"test"},
		},
	)

	assert.Nil(t, err)

	assert.Equal(t, user.ID, "69a495c9-7dbc-5733-e053-5b8c7c1155b0")
	assert.Equal(t, user.FirstName, "Oliver")
	assert.Equal(t, user.LastName, "Binns")
	assert.Equal(t, user.Username, "mail@oliverbinns.co.uk")

	// Visible App IDs are returned from the input as these are not available in the response:
	assert.Equal(t, user.VisibleAppIDs, []string{"test"})
}
