package appstore

import (
	"context"
	"net/http"

	"github.com/oliver-binns/appstore-go/authorization"
	"github.com/oliver-binns/appstore-go/users"
	"github.com/oliver-binns/googleplay-go/networking"
)

type Client struct {
	client  *networking.HTTPClient
	baseURL string
}

func AppStoreClient(
	keyID string,
	issuerID string,
	privateKey string,
) *Client {
	serviceAccount := authorization.Account{
		KeyID:      keyID,
		IssuerID:   issuerID,
		PrivateKey: privateKey,
	}

	tokenSource, err := authorization.NewTokenSource(serviceAccount)
	check(err)

	client := networking.NewAuthorizedClient(http.DefaultClient, tokenSource)

	return &Client{
		client:  &client,
		baseURL: "https://api.appstoreconnect.apple.com/v1/",
	}
}

func (c *Client) GetUser(ctx context.Context, id string) (*users.User, error) {
	return users.Get(*c.client, ctx, c.baseURL, id)
}

func (c *Client) CreateUser(ctx context.Context, user users.User) (*users.User, error) {
	return users.Create(*c.client, ctx, c.baseURL, user)
}

func (c *Client) ModifyUser(ctx context.Context, id string, user users.User) (*users.User, error) {
	return users.Modify(*c.client, ctx, c.baseURL, id, user)
}

func (c *Client) DeleteUser(ctx context.Context, id string) error {
	return users.Delete(*c.client, ctx, c.baseURL, id)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
