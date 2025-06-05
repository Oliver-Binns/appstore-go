package authorization

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenSource_Token(t *testing.T) {
	testStartTime := float64(time.Now().Unix())

	// Generate a token
	pk, b, err := generatePrivateKey()
	require.NoError(t, err)

	key := string(b[:])
	acc := mockAccount(key)

	source, err := NewTokenSource(acc)
	require.NoError(t, err)

	token, err := source.Token()
	require.NoError(t, err)

	// Parse the token
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return pk.Public(), nil
	})
	require.NoError(t, err)

	// The correct key was used:
	// Signed with EC
	assert.Equal(t, jwt.SigningMethodES256.Alg(), parsed.Method.Alg())
	assert.Equal(t, acc.KeyID, parsed.Header["kid"])

	claims, ok := parsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	// Contains the correct issuer and subject:
	assert.Equal(t, acc.IssuerID, claims["iss"])
	assert.Equal(t, "user", claims["sub"])

	// Contains the correct audience and scope:
	assert.Equal(t, "appstoreconnect-v1", claims["aud"])
	assert.Equal(t, "", claims["scope"])

	// Issued time is between the start of the test and now
	assert.GreaterOrEqual(t, claims["iat"], testStartTime)
	assert.LessOrEqual(t, claims["iat"], float64(time.Now().Unix()))

	// Expiration time is twenty minutes after issued time
	assert.Equal(t, claims["iat"].(float64)+1200, claims["exp"])
}

func TestTokenSource_TokenRefresh(t *testing.T) {
	// GIVEN I have a valid token
	_, b, err := generatePrivateKey()
	require.NoError(t, err)

	key := string(b[:])
	acc := mockAccount(key)
	source, err := ShortlivedTokenSource(acc)

	require.NoError(t, err)

	token, err := source.Token()
	require.NoError(t, err)

	same, err := source.Token()
	require.NoError(t, err)
	assert.Equal(t, token, same)

	// WHEN the token expires
	time.Sleep(time.Second)

	// THEN it should be refreshed: this currently fails!
	new, err := source.Token()
	require.NoError(t, err)
	assert.NotEqual(t, token, new)
}

func mockAccount(key string) Account {
	return Account{
		KeyID:      "2X9R4HXF34",
		IssuerID:   "57246542-96fe-1a63-e053-0824d011072a",
		PrivateKey: key,
	}
}

func generatePrivateKey() (*ecdsa.PrivateKey, []byte, error) {
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	keyBytes, err := x509.MarshalECPrivateKey(pk)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal private key: %w", err)
	}

	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	})

	return pk, pemBytes, nil
}

func ShortlivedTokenSource(account Account) (TokenSource, error) {
	pk, err := jwt.ParseECPrivateKeyFromPEM([]byte(account.PrivateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &tokenSource{
		account:   account,
		pk:        pk,
		expiresIn: time.Second,
	}, nil
}
