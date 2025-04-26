package routes

import (
	"github.com/Edouard127/lambda-api/api/models/request"
	"github.com/Edouard127/lambda-api/api/models/response"
	"github.com/Edouard127/lambda-api/internal"
	"github.com/gin-gonic/gin"
	"github.com/yeqown/memcached"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// LinkDiscord godoc
//
//	@Summary	Links a discord account to an existing auth token
//	@Tags		Authentication
//	@Accept		json
//	@Produce	json
//	@Param		login	body	request.DiscordLink		true	"Discord RPC oauth token"
//	@Success	200	{object}	response.Authentication		"Successfully linked the discord account"
//	@Failure	400	{object}	response.ValidationError	"Invalid or missing fields in the request"
//	@Failure	401	{object}	response.Error				"Invalid discord token"
//	@Failure	500	{object}	response.Error				"Internal server error"
//	@Router		/link/discord [post]
//	@Security 	Bearer
func LinkDiscord(ctx *gin.Context, cache memcached.Client) {
	logger := ctx.MustGet("logger").(*zap.Logger)
	player := ctx.MustGet("player").(response.Player)

	var link request.DiscordLink

	err := ctx.ShouldBindJSON(&link)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ValidationError{
			Message: "Required fields are missing or invalid",
			Errors:  err.Error(),
		})
		return
	}

	err = response.GetDiscord(link.Token, &player)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Error{
			Message: "Invalid discord token",
		})
		return
	}

	signed, err := internal.NewJwt(player)
	if err != nil {
		logger.Error("Error signing token", zap.Error(err))

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{
			Message: "Failed to create the token",
		})
		return
	}

	// Touch the party if the player owns one
	item, _ := cache.GetAndTouch(ctx.Request.Context(), 86400, player.Hash())
	if item != nil {
		cache.Touch(ctx.Request.Context(), string(item.Value), 86400)
	}

	ctx.AbortWithStatusJSON(http.StatusOK, response.Authentication{
		AccessToken: signed,
		ExpiresIn:   int64(time.Hour * 24),
		TokenType:   "Bearer",
	})
}
