package response

type Settings struct {
	// The maximum number of players in the party.
	// example: 10
	MaxPlayers int `json:"max_players" form:"max_players" xml:"max_players" binding:"required"`

	// Whether the party is public or not.
	// If false can only be joined by invite.
	// example: true
	Public bool `json:"public" form:"public" xml:"public" binding:"required"`

	// Whether the party can be listed or not.
	// example: true
	Listed bool `json:"listed" form:"listed" xml:"listed" binding:"required"`
}

var DefaultSettings = &Settings{
	MaxPlayers: 10,
	Public:     false,
	Listed:     false,
}
