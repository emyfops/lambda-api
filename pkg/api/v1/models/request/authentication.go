package request

type Authentication struct {
	// The player's username.
	// example: Notch
	Username string `json:"username" form:"username" xml:"username" binding:"required"`

	// The player's Mojang session hash.
	// example: 069a79f444e94726a5befca90e38aaf5
	Hash string `json:"hash" form:"hash" xml:"hash" binding:"required"`
}
