package authorization

import (
	"crypto/ecdsa"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Account struct {
	KeyID      string
	IssuerID   string
	PrivateKey string
}

type TokenSource interface {
	Token() (string, error)
}

func NewTokenSource(account Account) (TokenSource, error) {
	pk, err := jwt.ParseECPrivateKeyFromPEM([]byte(account.PrivateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &tokenSource{
		account:   account,
		pk:        pk,
		expiresIn: 20 * time.Minute,
	}, nil
}

type tokenSource struct {
	sync.Mutex

	account   Account
	pk        *ecdsa.PrivateKey
	expiresIn time.Duration
	bearer    string
	expireAt  time.Time
}

func (ts *tokenSource) Token() (string, error) {
	ts.Lock()
	defer ts.Unlock()

	if ts.isExpired() {
		return ts.refresh()
	}

	return ts.bearer, nil
}

func (ts *tokenSource) isExpired() bool {
	return time.Now().After(ts.expireAt)
}

func (ts *tokenSource) refresh() (string, error) {
	// Create JWT as defined in https://developer.apple.com/documentation/appstoreconnectapi/generating-tokens-for-api-requests
	iat := time.Now()
	exp := iat.Add(ts.expiresIn)

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss":   ts.account.IssuerID,
		"sub":   "user",
		"scope": "",
		"aud":   "appstoreconnect-v1",
		"iat":   iat.Unix(),
		"exp":   exp.Unix(),
	})
	token.Header["alg"] = "ES256"
	token.Header["typ"] = "JWT"
	token.Header["kid"] = ts.account.KeyID

	bearer, err := token.SignedString(ts.pk)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	ts.bearer = bearer
	ts.expireAt = exp

	return bearer, nil
}
