package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	// ErrTokenExpired means that the JWT is valid, but it has expired
	ErrTokenExpired = errors.New("token has expired")
	// ErrInvalidSignature means that the JWT has been tampered with
	ErrInvalidSignature = errors.New("invalid signature")

	standartHeader = &header{
		Algorith: "HS256",
		Type:     "JWT",
	}
)

type header struct {
	Algorith string `json:"alg"`
	Type     string `json:"typ"`
}

type payload struct {
	Username  string `json:"usr"`
	ExpiresAt int64  `json:"exp"`
}

// Authenticator handles the logic around generating
// and validating JWT tokens
type Authenticator struct {
	secret string
}

// NewAutheniticator creates new Authenticator from the given parameters.
func NewAutheniticator(secret string) *Authenticator {
	return &Authenticator{
		secret: secret,
	}
}

// NewTokenForUser generates a new JWT for the given username,
// signs it with a secret and returns it.
func (a *Authenticator) NewTokenForUser(username string) (string, error) {
	p := &payload{
		Username:  username,
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
	}

	base, err := tokenBase(standartHeader, p)
	if err != nil {
		return "", fmt.Errorf("error building token base: %w", err)
	}

	return fmt.Sprintf("%s.%s", base, a.hash(base)), nil
}

func (a *Authenticator) hash(src string) string {
	key := []byte(a.secret)
	h := hmac.New(sha256.New, key)
	if _, err := h.Write([]byte(src)); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// DecodeToken accepts a token,
// and if valid and unexpired returns the username of the user
// the token belongs to, otherwise it returns an error.
// If the token is expired, a ErrTokenExpired is returned.
// If the JWT has been tampered with, a ErrInvalidSignature is returned.
func (a *Authenticator) DecodeToken(token string) (string, error) {
	components := strings.Split(token, ".")
	if len(components) != 3 {
		return "", fmt.Errorf("invalid token format")
	}

	pDec, err := base64.StdEncoding.DecodeString(components[1])
	if err != nil {
		return "", fmt.Errorf("error decoding payload: %w", err)
	}
	var p payload
	if err := json.Unmarshal(pDec, &p); err != nil {
		return "", fmt.Errorf("error unmarshaling payload: %w", err)
	}

	if time.Now().Unix() > p.ExpiresAt {
		return "", ErrTokenExpired
	}

	base, err := tokenBase(standartHeader, &p)
	if err != nil {
		return "", fmt.Errorf("error building token base: %w", err)
	}

	if components[2] != a.hash(base) {
		return "", ErrInvalidSignature
	}

	return p.Username, nil
}

func tokenBase(header *header, payload *payload) (string, error) {
	hJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	hEnc := base64.StdEncoding.EncodeToString(hJSON)

	pJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	pEnc := base64.StdEncoding.EncodeToString(pJSON)

	return fmt.Sprintf("%s.%s", hEnc, pEnc), nil
}
