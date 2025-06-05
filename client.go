package appstore

import (
	"net/http"

	"github.com/oliver-binns/appstore-go/authorization"
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

func check(e error) {
	if e != nil {
		panic(e)
	}
}
