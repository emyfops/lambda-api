package endpoints

import (
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models"
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
// @Param token body string true "Discord identify token"
// @Param username body string true "Minecraft username"
// @Param hash body string true "Mojang session hash"
// @Error 401 {object} models.Error
// @Error 500 {object} models.Error
// @Success 200 {object} models.Authentication
// @Router /party/login [post]
func Login(ctx *gin.Context) {
	token := ctx.Param("token")       // Discord identify token
	username := ctx.Param("username") // Minecraft username
	hash := ctx.Param("hash")         // Mojang session hash

	// Check if the user is already connected
	player, err := models.GetPlayer(username, hash, token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "You either have an invalid account or the hash has expired, please reconnect to the server",
		})
		return
	}

	signed, err := auth.CreateJwtToken(player)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred while signing the token",
		})
	}

	ctx.AbortWithStatusJSON(http.StatusOK, models.Authentication{
		AccessToken: signed,
		ExpiresIn:   int64(time.Hour * 24),
		TokenType:   "Bearer",
	})
}
