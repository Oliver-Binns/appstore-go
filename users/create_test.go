package users

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser_MakesRequest(t *testing.T) {
	user := UserInvitation{
		FirstName: "Joseph",
		LastName:  "Bloggs",
		Email:     "joe.bloggs@example.com",
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
	assert.Equal(t, `{"data":{"type":"userInvitations","attributes":{"firstName":"Joseph","lastName":"Bloggs","email":"joe.bloggs@example.com","roles":["MARKETING"]}}}
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
					"email": "mail@oliverbinns.co.uk"
				}
			}
		}`,
	}

	user, _ := Create(
		httpClient, context.Background(), "https://example.com", UserInvitation{},
	)

	assert.Equal(t, user.FirstName, "Oliver")
	assert.Equal(t, user.LastName, "Binns")
	assert.Equal(t, user.Email, "mail@oliverbinns.co.uk")
}
