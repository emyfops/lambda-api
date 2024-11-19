package endpoints

import (
	"context"
	"github.com/Edouard127/lambda-api/internal/app/gonic"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/request"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
)

// EditParty 	godoc
//
//	@Summary	Edit a party
//	@Tags		Party
//	@Accept		json
//	@Produce	json
//	@Param		settings	body		request.Settings	false	"Party configuration"
//	@Success	202			{object}	response.Party
//	@Failure	400			{object}	response.Error
//	@Failure	403			{object}	response.ValidationError
//	@Failure	404			{object}	response.Error
//	@Router		/party/edit [patch]
//	@Security	ApiKeyAuth
func EditParty(ctx *gin.Context, client *redis.Client) {
	var party response.Party
	var settings request.Settings
	if err := ctx.Bind(&settings); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ValidationError{
			Message: "Required fields are missing or invalid",
			Errors:  err.Error(),
		})
		return
	}

	player := gonic.MustGet[response.Player](ctx, "player")

	err := client.HGetAll(context.Background(), player.String()).Scan(&party)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Error{
			Message: "You are not in a party",
		})
		return
	}

	if party.Leader != player {
		ctx.AbortWithStatusJSON(http.StatusForbidden, response.Error{
			Message: "You are not the leader of the party",
		})
		return
	}

	party.Settings = settings
	client.HSet(context.Background(), player.String(), "settings", settings)

	ctx.AbortWithStatusJSON(http.StatusAccepted, party)
}
