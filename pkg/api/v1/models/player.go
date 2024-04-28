package models

type Player struct {
	// The Minecraft player.
	Minecraft MinecraftPlayer `json:"player"`

	// The Discord user.
	Discord DiscordUser `json:"discord"`
}

// GetPlayer returns a new player with the given name, hash and token.
// Returns nil if the Minecraft or Discord account is invalid.
func GetPlayer(name, hash, token string) *Player {
	minecraft := GetMinecraft(name, hash)
	discord := GetDiscord(token)

	if minecraft == nil || discord == nil {
		return nil
	}

	return &Player{
		Minecraft: *minecraft,
		Discord:   *discord,
	}
}
