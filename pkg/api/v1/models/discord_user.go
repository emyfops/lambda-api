package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DiscordUser struct {
	// The user's ID.
	// example: 1234567890
	ID int `json:"id"`

	// The user's username.
	// example: John Doe
	Username string `json:"username"`

	// The user's name.
	// example: john.doe
	GlobalName string `json:"global_name"`
}

// GetDiscord authenticates a user with the Discord token.
func GetDiscord(token string) *DiscordUser {
	req, _ := http.NewRequest("GET", "https://discord.com/api/v9/users/@me", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var discordResp *DiscordUser
	_ = json.Unmarshal(body, &discordResp)

	return discordResp
}
