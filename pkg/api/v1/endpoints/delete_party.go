package endpoints

import (
	"encoding/json"
	"errors"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"github.com/yeqown/memcached"
	"go.uber.org/zap"
	"net/http"
)

// DeleteParty godoc
//
//	@Summary	Delete an existing party
//	@Tags		Party
//	@Accept		json
//	@Produce	json
//	@Success	200	"Ok"
//	@Failure	403	{object}	response.Error	"You are not the leader of the party"
//	@Failure	404	{object}	response.Error	"You are not in a party"
//	@Failure	500	{object}	response.Error	"Internal server error"
//	@Router		/party/delete [delete]
//	@Security 	Bearer
func DeleteParty(ctx *gin.Context, cache memcached.Client) {
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

	if party.Leader != player {
		ctx.AbortWithStatusJSON(http.StatusForbidden, response.Error{
			Message: "You are not the leader of the party",
		})
		return
	}

	flow.PublishAsync(party.JoinSecret, nil) // Throw an exception on the client to catch as null

	err = cache.Delete(ctx.Request.Context(), player.Hash())
	if errors.Is(err, memcached.ErrNotFound) {
		logger.Error("Tried to delete a party player map that doesn't exist", zap.Any("player", player), zap.Any("party", party), zap.Error(err))

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{
			Message: "Internal server error. Please try again later.",
		})
		return
	}

	err = cache.Delete(ctx.Request.Context(), party.JoinSecret)
	if errors.Is(err, memcached.ErrNotFound) {
		logger.Error("Tried to delete a party secret that doesn't exist", zap.Any("player", player), zap.Any("party", party), zap.Error(err))

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{
			Message: "Internal server error. Please try again later.",
		})
		return
	}

	partyCountTotal.WithLabelValues("v1").Dec()
	loggedInTotal.WithLabelValues("v1").Sub(float64(len(party.Players)))

	ctx.AbortWithStatus(http.StatusOK)
}
