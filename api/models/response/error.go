package response

type Error struct {
	// The error message.
	// example: Player not found
	Message string `json:"message"`
}

type ValidationError struct {
	// The error message.
	// example: Validation error
	Message string `json:"message"`

	// The error details.
	// example: "errors": "Key: 'Authentication.Token' Error:Field validation for 'Token' failed on the 'required' tag"
	Errors string `json:"errors"`
}
