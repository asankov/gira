package client

import "net/http"

// XAuthToken is the name of the header used for authentication
var XAuthToken = "x-auth-token"

// Client is the struct that is used to communicate
// with the games service.
type Client struct {
	addr       string
	httpClient *http.Client
}

// New returns a new client with the given address.
func New(addr string) (*Client, error) {
	return &Client{
		addr:       addr,
		httpClient: &http.Client{},
	}, nil
}
