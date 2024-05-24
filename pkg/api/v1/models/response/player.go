package response

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Player struct {
	// The player's name.
	// example: Notch
	Name string `json:"name"`

	// The player's UUID.
	// example: 069a79f4-44e9-4726-a5be-fca90e38aaf5
	UUID string `json:"id"`

	// The player's Discord ID.
	// example: "385441179069579265"
	DiscordID string `json:"discord_id"`
}

func (pl *Player) String() string {
	return fmt.Sprintf("Player{Name: %s, UUID: %s, DiscordID: %s}", pl.Name, pl.UUID, pl.DiscordID)
}

// Because both the Minecraft and Discord API return the same structure,
// we can use a shared structure to unmarshal the response.
type sharedPlayer struct {
	Name string `json:"name"` // mojang only
	ID   string `json:"id"`   // discord and mojang
}

// GetPlayer returns a new player with the given name, hash and token.
// Returns nil if the Minecraft or Discord account is invalid.
func GetPlayer(token, name, hash string) (pl Player, err error) {
	err = GetMinecraft(name, hash, &pl)
	if err != nil {
		return
	}

	err = GetDiscord(token, &pl)

	return
}

// GetMinecraft authenticates a user with the Mojang session hash.
// It is used to prove that a user owns a Minecraft account and is connected
// to a server without requiring OAuth2 authentication.
// However, the session hash is only valid for a very small time-frame.
func GetMinecraft(name, hash string, player *Player) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://sessionserver.mojang.com/session/minecraft/hasJoined?username=%s&serverId=%s", name, hash), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var shared sharedPlayer

	err = json.Unmarshal(body, &shared)
	if err != nil {
		return err
	}

	player.Name = shared.Name
	player.UUID = shared.ID

	return resp.Body.Close()
}

// GetDiscord authenticates a user with the Discord token.
func GetDiscord(token string, player *Player) error {
	req, _ := http.NewRequest("GET", "https://discord.com/api/v10/users/@me", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var shared sharedPlayer

	err = json.Unmarshal(body, &shared)
	if err != nil {
		return err
	}

	player.DiscordID = shared.ID

	return resp.Body.Close()
}
