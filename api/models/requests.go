package models

import "github.com/google/uuid"

type Authentication struct {
	// The player's username.
	// example: Notch
	Username string `json:"username" binding:"required"`

	// The player's Mojang session hash.
	// example: 069a79f444e94726a5befca90e38aaf5
	Hash string `json:"hash" binding:"required"`
}

type DiscordLink struct {
	// The player's authentication Discord token
	// 	example: OTk1MTU1NzcyMzYxMTQ2NDM4
	Token string `json:"token" form:"token" xml:"token" binding:"required"`
}

type JoinParty struct {
	// The join secret of the party.
	// 	example: "RzTMeBHZu3VoNEUpNQFSnMYpSNQgWQ2rYM4u3RHSPQIacMxE4KH63OYQEDAD0P0bnjZBYUFHB0I5wOx2xEiXOU9SKBxRZ5YcvYjjiZWhsGGss3vnRajvFn4trJgTnvz2"
	Secret string `json:"secret" binding:"required"`
}

type CapeLookup struct {
	Players []uuid.UUID `json:"players" binding:"required"`
}
