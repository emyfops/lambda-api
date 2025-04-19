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

// LeaveParty allows a player to leave a party.
//
//	@Summary	Leave a party
//	@Tags		Party
//	@Accept		json
//	@Produce	json
//	@Success	200	"Ok"
//	@Failure	404	{object}	response.Error	"You are not in a party"
//	@Failure	500	{object}	response.Error	"Internal server error"
//	@Router		/party/leave [put]
//	@Security 	Bearer
func LeaveParty(ctx *gin.Context, cache memcached.Client) {
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

	if player == party.Leader {
		DeleteParty(ctx, cache)
		return
	}

	// Close the channel if the player has subscribed to the party events
	channel, ok := subscriptions[player]
	if ok {
		close(channel)
	}

	cache.Delete(ctx.Request.Context(), player.Hash())

	party.Remove(player)
	party.Update(cache)

	flow.PublishAsync(party.JoinSecret, party)

	ctx.AbortWithStatus(http.StatusOK)
}
