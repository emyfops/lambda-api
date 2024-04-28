package endpoints

import (
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models/request"
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models/response"
	"github.com/Edouard127/lambda-rpc/pkg/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Login godoc
// @BasePath /api/v1
// @Summary Login to the server
// @Description Login to the server using a Discord identify token, a Minecraft username and a Mojang session hash
// @Tags Party
// @Accept json
// @Produce json
// @Param login body request.Authentication true "Authentication"
// @Success 200 {object} response.Authentication
// @Router /party/login [post]
func Login(ctx *gin.Context) {
	var login request.Authentication
	if err := ctx.ShouldBind(&login); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	// Check if the user is already connected
	player, err := response.GetPlayer(login.Token, login.Username, login.Hash)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "You either have an invalid discord account or the hash has expired, please reconnect to the server",
		})
		return
	}

	signed, err := auth.CreateJwtToken(player)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred while signing the token",
		})
	}

	ctx.AbortWithStatusJSON(http.StatusOK, response.Authentication{
		AccessToken: signed,
		ExpiresIn:   int64(time.Hour * 24),
		TokenType:   "Bearer",
	})
}
