package middlewares

import (
	"github.com/Edouard127/lambda-api/api/models/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// DiscordCheck sits between the authorization middleware and the handler function
// and ensure that the player has linked their discord account
func DiscordCheck(ctx *gin.Context) {
	player := ctx.MustGet("player").(response.Player)

	if !player.HasDiscord() {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized,
			gin.H{"error": "You did not link your discord account"})
		return
	}

	ctx.Next()
}
