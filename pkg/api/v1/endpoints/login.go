package endpoints

import (
	"github.com/Edouard127/lambda-rpc/internal/app/auth"
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models/request"
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Login godoc
// @BasePath /api/v1
// @Summary Login to the server
// @Description Login to the server using a Discord identify token, a Minecraft username and a Mojang session hash
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login body request.Authentication true "Authentication"
// @Success 200 {object} response.Authentication
// @Failure 400 {object} response.ValidationError
// @Failure 401 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /party/login [post]
func Login(ctx *gin.Context) {
	var login request.Authentication
	if err := ctx.Bind(&login); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ValidationError{
			Message: "Required fields are missing or invalid",
			Errors:  err.Error(),
		})
		return
	}

	player, err := response.GetPlayer(login.Token, login.Username, login.Hash)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Error{
			Message: "Invalid credentials",
		})
		return
	}

	signed, err := auth.CreateJwtToken(player)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{
			Message: "Failed to create token",
		})
	}

	ctx.AbortWithStatusJSON(http.StatusOK, response.Authentication{
		AccessToken: signed,
		ExpiresIn:   int64(time.Hour * 24),
		TokenType:   "Bearer",
	})
}
