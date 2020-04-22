package auth

import (
	"testing"
)

var (
	expectedUsername = "username"
)

func TestToken(t *testing.T) {

	a := NewAutheniticator("secret")

	token, err := a.NewTokenForUser(expectedUsername)
	if err != nil {
		t.Fatalf("got (%v), expected nil error when creating token for user", err)
	}

	username, err := a.DecodeToken(token)
	if err != nil {
		t.Fatalf("got (%v), expected nil error when decoding token", err)
	}

	if username != expectedUsername {
		t.Errorf("got (%v), expected (%v) for username", username, expectedUsername)
	}
}
