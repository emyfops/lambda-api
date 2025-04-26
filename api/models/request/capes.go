package request

import "github.com/google/uuid"

type CapeLookup struct {
	Players []uuid.UUID `json:"players" binding:"required"`
}
