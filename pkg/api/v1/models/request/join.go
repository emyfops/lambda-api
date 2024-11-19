package request

import "github.com/google/uuid"

type JoinParty struct {
	// The ID of the party.
	// example: "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
	ID uuid.UUID `json:"id" form:"id" xml:"id" binding:"required"`
}
