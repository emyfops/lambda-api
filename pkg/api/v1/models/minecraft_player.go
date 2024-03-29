package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MinecraftPlayer struct {
	// The player's name.
	// example: Notch
	Name string `json:"name"`

	// The player's UUID.
	// example: 069a79f4-44e9-4726-a5be-fca90e38aaf5
	UUID string `json:"id"`
}

// GetMinecraft authenticates a user with the Mojang session hash.
// It is used to prove that a user owns a Minecraft account and is connected
// to a server without requiring OAuth2 authentication.
// However, the session hash is only valid for a very small time-frame.
func GetMinecraft(name, hash string) *MinecraftPlayer {
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://sessionserver.mojang.com/session/minecraft/hasJoined?username=%s&serverId=%s", name, hash), nil)

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var mojangResp *MinecraftPlayer
	_ = json.Unmarshal(body, &mojangResp)

	return mojangResp
}
