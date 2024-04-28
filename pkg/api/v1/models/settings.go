package models

type Settings struct {
	// The maximum number of players in the party.
	// example: 10
	MaxPlayers int `json:"max_players"`

	// Whether the party is public or not.
	// If false can only be joined by invite.
	// example: true
	Public bool `json:"public"`

	// Whether the party can be listed or not.
	// example: true
	Listed bool `json:"listed"`
}

var DefaultSettings = &Settings{
	MaxPlayers: 10,
	Public:     false,
	Listed:     false,
}
