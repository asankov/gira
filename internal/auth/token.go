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

	"github.com/gira-games/api/pkg/models"
)

var (
	// ErrTokenExpired means that the JWT is valid, but it has expired
	ErrTokenExpired = errors.New("token has expired")
	// ErrInvalidSignature means that the JWT has been tampered with
	ErrInvalidSignature = errors.New("invalid signature")
	// ErrInvalidFormat means that the token is not a valid JWT token
	ErrInvalidFormat = errors.New("invalid token format")

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
	User      *models.User
	ExpiresAt int64 `json:"exp"`
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
// with the default expiration of 50 minutes, signs it with a secret and returns it.
func (a *Authenticator) NewTokenForUser(user *models.User) (string, error) {
	return a.NewTokenForUserWithExpiration(user, 50*time.Minute)
}

// NewTokenForUserWithExpiration generates a new JWT for the given username,
// with expiration now + d, signs it with a secret and returns it.
func (a *Authenticator) NewTokenForUserWithExpiration(user *models.User, d time.Duration) (string, error) {
	p := &payload{
		User:      user,
		ExpiresAt: time.Now().Add(d).Unix(),
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
func (a *Authenticator) DecodeToken(token string) (*models.User, error) {
	components := strings.Split(token, ".")
	if len(components) != 3 {
		return nil, ErrInvalidFormat
	}

	pDec, err := base64.StdEncoding.DecodeString(components[1])
	if err != nil {
		return nil, fmt.Errorf("error decoding payload: %w", err)
	}
	var p payload
	if err := json.Unmarshal(pDec, &p); err != nil {
		return nil, fmt.Errorf("error unmarshaling payload: %w", err)
	}

	if time.Now().Unix() > p.ExpiresAt {
		return nil, ErrTokenExpired
	}

	base, err := tokenBase(standartHeader, &p)
	if err != nil {
		return nil, fmt.Errorf("error building token base: %w", err)
	}

	if components[2] != a.hash(base) {
		return nil, ErrInvalidSignature
	}

	return p.User, nil
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
