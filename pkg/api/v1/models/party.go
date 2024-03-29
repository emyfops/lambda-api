package models

import (
	"github.com/Edouard127/lambda-rpc/internal/util"
	"time"
)

// The Party represents a Discord party.
// It contains an ID, a creation date and a list of players.
// There is no leader in the party, all players have the same rights.
type Party struct {
	// The ID of the party.
	// It is a random string of 69 characters.
	ID string `json:"id"`

	// The creation date of the party.
	// example: 2021-10-10T12:00:00Z
	Creation time.Time `json:"creation"`

	// The list of players in the party.
	Players []Player `json:"players"`
}

// New creates a new party with the given name and players.
// The UUID is generated automatically.
func New(players ...Player) *Party {
	return &Party{
		ID:       util.RandStringBytesMaskSrcUnsafe(69),
		Creation: time.Now(),
		Players:  players,
	}
}

func (pt *Party) Add(player Player) {
	pt.Players = append(pt.Players, player)
}

func (pt *Party) Remove(player Player) {
	for i, p := range pt.Players {
		if p.Minecraft.UUID == player.Minecraft.UUID {
			pt.Players = append(pt.Players[:i], pt.Players[i+1:]...)
			return
		}
	}
}
