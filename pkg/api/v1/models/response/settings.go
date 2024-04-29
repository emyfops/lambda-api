package response

type Settings struct {
	// The maximum number of players in the party.
	// example: 10
	MaxPlayers int `json:"max_players" form:"max_players" xml:"max_players" binding:"required"`

	// Whether the party is public or not.
	// If false can only be joined by Discord invites.
	// If true can be joined by anyone with the join secret
	// example: true
	Public bool `json:"public" form:"public" xml:"public" binding:"required"`
}

var DefaultSettings = &Settings{
	MaxPlayers: 10,
	Public:     false,
}
