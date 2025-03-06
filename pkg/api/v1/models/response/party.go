package response

import (
	"context"
	"encoding/json"
	"github.com/Edouard127/lambda-api/internal"
	"github.com/google/uuid"
	"github.com/yeqown/memcached"
	"time"
)

// Party represents a lobby of players
type Party struct {
	// The ID of the party
	// 	It's a random UUID
	ID uuid.UUID `json:"id" redis:"id"`

	// The join secret of the party
	// 	example: "RzTMeBHZu3VoNEUpNQFSnMYpSNQgWQ2rYM4u3RHSPQIacMxE4KH63OYQEDAD0P0bnjZBYUFHB0I5wOx2xEiXOU9SKBxRZ5YcvYjjiZWhsGGss3vnRajvFn4trJgTnvz2"
	JoinSecret string `json:"join_secret"`

	// The leader of the party
	Leader Player `json:"leader"`

	// The creation date of the party
	// 	example: 2021-10-10T12:00:00Z
	Creation time.Time `json:"creation"`

	// The list of players in the party
	Players []Player `json:"players"`
}

// NewParty returns a new party with the given leader
func NewParty(leader Player) Party {
	return Party{
		ID:         uuid.New(),
		JoinSecret: internal.RandString(128),
		Leader:     leader,
		Creation:   time.Now(),
		Players:    []Player{leader},
	}
}

func (pt *Party) Add(player Player) {
	pt.Players = append(pt.Players, player)
}

func (pt *Party) Remove(player Player) {
	for i, p := range pt.Players {
		if p.UUID == player.UUID {
			pt.Players = append(pt.Players[:i], pt.Players[i+1:]...)
			return
		}
	}
}

func (pt *Party) Update(cache memcached.Client) {
	bytes, _ := json.Marshal(pt)

	for _, player := range pt.Players {
		cache.Set(context.Background(), player.Hash(), bytes, 0, 0)
	}
}
