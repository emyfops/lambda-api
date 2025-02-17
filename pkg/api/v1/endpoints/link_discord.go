package endpoints

import (
	"github.com/Edouard127/lambda-api/internal"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/request"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// LinkDiscord godoc
//
//	@Summary 	Links a discord account to an existing auth token
//	@Tags		Authentication
//	@Accept 	json
//	@Produce	json
//	@Param		login	body		request.DiscordLink		true
//	@Success	200		{object}	response.Authentication
//	@Failure	400		{object}	response.ValidationError
//	@Failure	401		{object}	response.Error
//	@Failure	500		{object}	response.Error
//	@Router		/link/discord 		[post]
func LinkDiscord(ctx *gin.Context) {
	var link request.DiscordLink
	if err := ctx.Bind(&link); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ValidationError{
			Message: "Required fields are missing or invalid",
			Errors:  err.Error(),
		})
		return
	}

	player := ctx.MustGet("player").(response.Player)
	err := response.GetDiscord(link.Token, &player)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Error{
			Message: "Invalid discord token",
		})
		return
	}

	signed, err := internal.NewJwt(player)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{
			Message: "Failed to create token",
		})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, response.Authentication{
		AccessToken: signed,
		ExpiresIn:   int64(time.Hour * 24),
		TokenType:   "Bearer",
	})
}
