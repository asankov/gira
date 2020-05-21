package auth

import (
	"testing"

	"github.com/asankov/gira/pkg/models"
)

var (
	expectedUser     = &models.User{Username: expectedUsername}
	expectedUsername = "username"
)

func TestToken(t *testing.T) {

	a := NewAutheniticator("secret")

	token, err := a.NewTokenForUser(expectedUser)
	if err != nil {
		t.Fatalf("got (%v), expected nil error when creating token for user", err)
	}

	usr, err := a.DecodeToken(token)
	if err != nil {
		t.Fatalf("got (%v), expected nil error when decoding token", err)
	}

	got, actual := usr.Username, expectedUser.Username
	if got != actual {
		t.Errorf("got (%v), expected (%v) for username", got, actual)
	}
}
