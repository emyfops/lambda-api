package routes

import (
	"encoding/json"
	"errors"
	"github.com/Edouard127/lambda-api/api/models/response"
	"github.com/gin-gonic/gin"
	"github.com/yeqown/memcached"
	"go.uber.org/zap"
	"net/http"
)

// GetParty godoc
//
//	@Summary	Get the party of the player
//	@Tags		Party
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	response.Party	"Details of the player's current party"
//	@Failure	404	{object}	response.Error	"You are not in a party"
//	@Failure	500	{object}	response.Error	"Internal server error"
//	@Router		/party [get]
//	@Security 	Bearer
func GetParty(ctx *gin.Context, cache memcached.Client) {
	logger := ctx.MustGet("logger").(*zap.Logger)
	player := ctx.MustGet("player").(response.Player)

	item, err := cache.Get(ctx.Request.Context(), player.Hash())
	if !errors.Is(err, memcached.ErrNotFound) && err != nil {
		logger.Error("Error getting party from cache", zap.Any("player", player), zap.Error(err))

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{
			Message: "Internal server error. Please try again later.",
		})
		return
	}
	if errors.Is(err, memcached.ErrNotFound) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Error{
			Message: "You are not in a party",
		})
		return
	}

	var party response.Party
	json.Unmarshal(item.Value, &party)

	ctx.AbortWithStatusJSON(http.StatusOK, party)
}
