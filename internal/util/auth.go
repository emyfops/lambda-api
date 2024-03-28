package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Player struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

// IsConnected authenticates a user with the Mojang session hash.
// It is used to prove that a user owns a Minecraft account and is connected
// to a server without requiring OAuth2 authentication.
// However, the session hash is only valid for a very small time-frame.
func IsConnected(name, hash string) *Player {
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://sessionserver.mojang.com/session/minecraft/hasJoined?username=%s&serverId=%s", name, hash), nil)

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var mojangResp *Player
	json.Unmarshal(body, &mojangResp)

	return mojangResp
}
