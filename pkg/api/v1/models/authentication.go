package models

// Authentication represents the response from the authentication endpoint.
type Authentication struct {
	// The access token to use for the API
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
	AccessToken string `json:"access_token"`

	// The duration of the token (in seconds).
	// example: 3600
	ExpiresIn int64 `json:"expires_in"`

	// The type of the token.
	// example: Bearer
	TokenType string `json:"token_type"`
}
