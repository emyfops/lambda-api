package routes

import (
	"encoding/json"
	"errors"
	"github.com/Edouard127/lambda-api/api/metrics"
	"github.com/Edouard127/lambda-api/api/models/response"
	"github.com/gin-gonic/gin"
	"github.com/yeqown/memcached"
	"go.uber.org/zap"
	"net/http"
)

// CreateParty godoc
//
//	@Summary	Create a new party
//	@Tags		Party
//	@Accept		json
//	@Produce	json
//	@Success	201			{object}	response.Party	"Successfully created the party"
//	@Failure	400			{object}	response.ValidationError	"Invalid request parameters"
//	@Failure	409			{object}	response.Error	"Conflict - Player is already in a party"
//	@Failure	500			{object}	response.Error	"Internal server error"
//	@Router		/party/create [post]
//	@Security 	Bearer
func CreateParty(ctx *gin.Context, cache memcached.Client) {
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
	if !errors.Is(err, memcached.ErrNotFound) {
		var party response.Party
		json.Unmarshal(item.Value, &party)

		if party.Leader == player {
			DeleteParty(ctx, cache)
		} else {
			LeaveParty(ctx, cache)
		}

		return
	}

	party := response.NewParty(player)
	bytes, _ := json.Marshal(party)

	err = cache.Set(ctx.Request.Context(), player.Hash(), bytes, 0, 86400)
	if err != nil {
		logger.Error("Error mapping player to party", zap.Error(err))
		goto err500
	}

	err = cache.Set(ctx.Request.Context(), party.JoinSecret, bytes, 0, 86400)
	if err != nil {
		logger.Error("Error mapping join secret to party", zap.Error(err))
		goto err500
	}

	metrics.PartyCountTotal.Inc()

	ctx.AbortWithStatusJSON(http.StatusCreated, party)
	return

err500:
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{
		Message: "Internal server error. Please try again later.",
	})
	return
}
