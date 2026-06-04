package appstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseURL(t *testing.T) {
	key := `-----BEGIN PRIVATE KEY-----
MHcCAQEEIG706QZ+qBP9FxNbs8lVhIf0w/hJJ+pMu6YtG/d8uqnkoAoGCCqGSM49
AwEHoUQDQgAEnMKTGhM0U4Q5rCvgobWZQtcmknAEZOxTqjmtJf1jlTfHO7iLykAj
AoyVWzvsnOZ2F3ujWssdv6b27lkdrm513w==
-----END PRIVATE KEY-----`

	client := AppStoreClient(
		"key-id",
		"issuer-id",
		key,
	)

	assert.Equal(t, "https://api.appstoreconnect.apple.com/v1/", client.baseURL)
}
