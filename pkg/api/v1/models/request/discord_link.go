package request

type DiscordLink struct {
	// The player's authentication Discord token
	// 	example: OTk1MTU1NzcyMzYxMTQ2NDM4
	Token string `json:"token" form:"token" xml:"token" binding:"required"`
}
