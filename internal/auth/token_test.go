package auth_test

import (
	"testing"
	"time"

	"github.com/asankov/gira/internal/auth"
	"github.com/asankov/gira/internal/fixtures/assert"
	"github.com/asankov/gira/pkg/models"
)

var (
	expectedUser     = &models.User{Username: expectedUsername}
	expectedUsername = "username"

	authenticator = auth.NewAutheniticator("secret")
)

func TestToken(t *testing.T) {
	token, err := authenticator.NewTokenForUser(expectedUser)
	if err != nil {
		t.Fatalf("got (%v), expected nil error when creating token for user", err)
	}

	usr, err := authenticator.DecodeToken(token)
	if err != nil {
		t.Fatalf("got (%v), expected nil error when decoding token", err)
	}

	got, actual := usr.Username, expectedUser.Username
	if got != actual {
		t.Errorf("got (%v), expected (%v) for username", got, actual)
	}
}

func TestDecodeTokenError(t *testing.T) {
	_, err := authenticator.DecodeToken("o.o")

	assert.Error(t, err, auth.ErrInvalidFormat)
}

func TestTokenExpired(t *testing.T) {
	token, err := authenticator.NewTokenForUserWithExpiration(expectedUser, 1*time.Millisecond)
	if err != nil {
		t.Fatalf("got (%v), expected nil error when creating token for user", err)
	}

	// wait for the token to expire
	time.Sleep(1 * time.Second)

	_, err = authenticator.DecodeToken(token)
	assert.Error(t, err, auth.ErrTokenExpired)
}

func TestInvalidSignature(t *testing.T) {
	newAuthenticator := auth.NewAutheniticator("some.other.secret")
	newToken, err := newAuthenticator.NewTokenForUser(&models.User{})
	assert.NoError(t, err)

	_, err = authenticator.DecodeToken(newToken)
	assert.Error(t, err, auth.ErrInvalidSignature)
}
