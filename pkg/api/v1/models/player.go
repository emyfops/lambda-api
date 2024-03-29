package models

type Player struct {
	// The Minecraft player.
	Minecraft MinecraftPlayer `json:"player"`

	// The Discord user.
	Discord DiscordUser `json:"discord"`
}

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
