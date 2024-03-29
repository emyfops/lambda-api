package models

type Player struct {
	Minecraft MinecraftPlayer `json:"player"`
	Discord   DiscordUser     `json:"discord"`
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
