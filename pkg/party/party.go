package party

import (
	"github.com/Edouard127/lambda-rpc/internal/util"
	"time"
)

type Party struct {
	Name     string        `json:"name"`
	ID       string        `json:"id"`
	Creation time.Time     `json:"creation"`
	Players  []util.Player `json:"players"`
}

// New creates a new party with the given name and players.
// The ID is generated automatically.
func New(name string, players ...util.Player) *Party {
	return &Party{
		Name:     name,
		ID:       util.RandStringBytesMaskSrcUnsafe(69),
		Creation: time.Now(),
		Players:  players,
	}
}

func (pt *Party) Add(player util.Player) {
	pt.Players = append(pt.Players, player)
}

func (pt *Party) Remove(player util.Player) {
	for i, p := range pt.Players {
		if p.ID == player.ID {
			pt.Players = append(pt.Players[:i], pt.Players[i+1:]...)
			return
		}
	}
}
