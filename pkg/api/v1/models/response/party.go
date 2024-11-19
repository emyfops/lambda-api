package response

import (
	"github.com/Edouard127/lambda-api/internal/app/random"
	"time"
)

// The Party represents a Discord party.
// It contains an ID, a creation date and a list of players.
type Party struct {
	// The ID of the party.
	// It is a random string of 30 characters.
	ID string `json:"id"`

	// The join secret of the party.
	// It is a random string of 100 characters.
	JoinSecret string `json:"join_secret"`

	// The leader of the party.
	Leader Player `json:"leader"`

	// The creation date of the party.
	// example: 2021-10-10T12:00:00Z
	Creation time.Time `json:"creation"`

	// The list of players in the party.
	Players []Player `json:"players"`

	// The settings of the party.
	Settings Settings `json:"settings"`
}

// NewPartyWithSettings returns a new party with the given leader and settings.
func NewPartyWithSettings(leader Player, settings *Settings) *Party {
	if settings == nil {
		settings = DefaultSettings
	}

	return &Party{
		ID:         random.RandString(30),
		JoinSecret: random.RandString(100),
		Leader:     leader,
		Creation:   time.Now(),
		Players:    []Player{leader},
		Settings:   *settings,
	}
}

// NewParty returns a new party with the given leader and default settings.
func NewParty(leader Player) *Party {
	return NewPartyWithSettings(leader, DefaultSettings)
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
