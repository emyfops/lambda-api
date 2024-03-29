package models

import "time"

// Authentication represents the response from the authentication endpoint.
type Authentication struct {
	// The access token to use for the API
	AccessToken string `json:"access_token"`

	// The duration of the token (in seconds).
	ExpiresIn time.Duration `json:"expires_in"`

	// The type of the token.
	TokenType string `json:"token_type"`

	// The message to display to the user.
	Message string `json:"message"`
}
