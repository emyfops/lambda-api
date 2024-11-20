package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Edouard127/lambda-api/pkg/cmd"
	"io"
	"net/http"
)

var ErrCouldNotVerifyMinecraft = errors.New("could not verify minecraft account")

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

	// Whether the player is marked as unsafe.
	// example: true
	Unsafe bool `json:"unsafe"`
}

func (pl *Player) String() string {
	return fmt.Sprintf("Player{Name: %s, UUID: %s, DiscordID: %s, Unsafe: %t}", pl.Name, pl.UUID, pl.DiscordID, pl.Unsafe)
}

// Because both the Minecraft and Discord API return the same structure,
// we can use a shared structure to unmarshal the response.
type sharedPlayer struct {
	Name string `json:"name"` // mojang only
	ID   string `json:"id"`   // discord and mojang
}

// GetPlayer returns a new player with the given name, hash and token.
func GetPlayer(token, name, hash string) (pl Player, err error) {
	err = GetMinecraft(name, hash, &pl)
	if errors.Is(err, ErrCouldNotVerifyMinecraft) &&
		cmd.Arguments().AllowInsecure {
		// If the Minecraft account is invalid, we can still try to authenticate the player with Discord.
		pl.Unsafe = true
		err = nil
	} else {
		return
	}

	err = GetDiscord(token, &pl)
	return
}

// GetMinecraft authenticates a user with the Mojang session hash.
//
//		This function proves the authenticity of a Minecraft username by checking
//		the session hash provided by the client.  The hash originally comes
//		from the server and is generated using the Yggdrasil Public Key.
//	 This is not the most secure way of authenticating users since the hash
//	 is susceptible to replay attacks.  However, it is the only way to
//	 authenticate users without requesting access to their Microsoft account.
func GetMinecraft(name, hash string, player *Player) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://sessionserver.mojang.com/session/minecraft/hasJoined?username=%s&serverId=%s", name, hash), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)

	if resp.StatusCode == http.StatusNoContent {
		// We can assume that either the hash or the username is invalid.
		return ErrCouldNotVerifyMinecraft
	}

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
