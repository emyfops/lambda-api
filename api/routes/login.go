package routes

import (
	"github.com/Edouard127/lambda-api/api/metrics"
	"github.com/Edouard127/lambda-api/api/models/request"
	"github.com/Edouard127/lambda-api/api/models/response"
	"github.com/Edouard127/lambda-api/internal"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// Login allows a player to log in to the server using a Minecraft username and Mojang session hash.
//
//	@Summary	Login to the server
//	@Tags		Authentication
//	@Accept		json
//	@Produce	json
//	@Param		login	body	request.Authentication	true	"Authentication credentials (Minecraft username and Mojang session hash)"
//	@Success	200	{object}	response.Authentication		"Successfully logged in and retrieved authentication token"
//	@Failure	400	{object}	response.ValidationError	"Invalid or missing authentication fields"
//	@Failure	401	{object}	response.Error				"Invalid credentials"
//	@Failure	500	{object}	response.Error				"Internal server error"
//	@Router		/login [post]
func Login(ctx *gin.Context) {
	logger := ctx.MustGet("logger").(*zap.Logger)

	var login request.Authentication

	err := ctx.ShouldBindJSON(&login)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ValidationError{
			Message: "Required fields are missing or invalid",
			Errors:  err.Error(),
		})
		return
	}

	player, err := response.GetPlayer(login.Username, login.Hash)
	if err != nil {
		metrics.FailedLogins.Inc()
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Error{
			Message: "Invalid credentials",
		})
		return
	}

	signed, err := internal.NewJwt(player)
	if err != nil {
		logger.Error("Error signing token", zap.Error(err))
		metrics.FailedLogins.Inc()

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{
			Message: "Failed to create token",
		})
		return
	}

	metrics.SuccessfulLogins.Inc()

	ctx.AbortWithStatusJSON(http.StatusOK, response.Authentication{
		AccessToken: signed,
		ExpiresIn:   int64(time.Hour * 24),
		TokenType:   "Bearer",
	})
}
