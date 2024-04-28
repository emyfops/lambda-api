package models

// The Error represents an error.
// It contains a message.
type Error struct {
	// The message of the error.
	// example: "You either have an invalid account or the hash has expired, please reconnect to the server"
	Message string `json:"message"`

	// The code of the error.
	// example: 401
	Code int `json:"code"`

	// The status of the error.
	// example: "Unauthorized"
	Status string `json:"status"`

	// The timestamp of the error.
	// example: 2021-10-10T12:00:00Z
	Timestamp string `json:"timestamp"`

	// The path of the error.
	// example: "/api/v1/party/login"
	Path string `json:"path"`

	// The method of the error.
	// example: "POST"
	Method string `json:"method"`
}
