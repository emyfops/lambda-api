package request

type Authentication struct {
	// The player's Discord token.
	// example: OTk1MTU1NzcyMzYxMTQ2NDM4.AAAAAA.BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB
	Token string `json:"token" form:"token" xml:"token" binding:"required"`

	// The player's username.
	// example: Notch
	Username string `json:"username" form:"username" xml:"username" binding:"required"`

	// The player's Mojang session hash.
	// example: 069a79f444e94726a5befca90e38aaf5
	Hash string `json:"hash" form:"hash" xml:"hash" binding:"required"`
}
