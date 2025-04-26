package response

import "github.com/google/uuid"

type Cape struct {
	Uuid uuid.UUID `json:"uuid"`
	Type string    `json:"type"`
}
