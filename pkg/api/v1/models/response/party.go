package response

import (
	"github.com/Edouard127/lambda-api/internal/app/random"
	"github.com/google/uuid"
	"time"
)

// Party represents a lobby of players
type Party struct {
	// The ID of the party.
	// It is a random UUID.
	ID uuid.UUID `json:"id" redis:"id"`

	// The join secret of the party.
	// It is a random string of 100 characters.
	JoinSecret string `json:"join_secret" redis:"join_secret"`

	// The leader of the party.
	Leader Player `json:"leader" redis:"leader"`

	// The creation date of the party.
	// example: 2021-10-10T12:00:00Z
	Creation time.Time `json:"creation" redis:"creation"`

	// The list of players in the party.
	Players []Player `json:"players" redis:"players"`

	// The settings of the party.
	Settings Settings `json:"settings" redis:"settings"`
}

// NewParty returns a new party with the given leader and settings.
func NewParty(leader Player, settings *Settings) *Party {
	if settings == nil {
		settings = DefaultSettings
	}

	return &Party{
		ID:         uuid.New(),
		JoinSecret: random.RandString(100),
		Leader:     leader,
		Creation:   time.Now(),
		Players:    []Player{leader},
		Settings:   *settings,
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
