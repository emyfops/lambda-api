package endpoints

import (
	"encoding/json"
	"errors"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/request"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/yeqown/memcached"
	"go.uber.org/zap"
	"net/http"
)

var loggedInTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "lambda_api_logged_in_users",
	Help: "Total number of logged in users",
}, []string{"version"})

// JoinParty godoc
//
//	@Summary	Join a party
//	@Tags		Party
//	@Accept		json
//	@Produce	json
//	@Param		id	body		request.JoinParty	true	"Party secret"
//	@Success	202	{object}	response.Party				"Successfully joined the party"
//	@Failure	400	{object}	response.ValidationError	"Invalid request, required fields are missing or incorrect"
//	@Failure	404	{object}	response.Error				"The party does not exist"
//	@Failure	500	{object}	response.Error				"Internal server error"
//	@Router		/party/join [put]
//	@Security 	Bearer
func JoinParty(ctx *gin.Context, cache memcached.Client) {
	logger := ctx.MustGet("logger").(*zap.Logger)
	player := ctx.MustGet("player").(response.Player)

	var join request.JoinParty

	err := ctx.Bind(&join)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ValidationError{
			Message: "Required fields are missing or invalid",
			Errors:  err.Error(),
		})
		return
	}

	item, err := cache.Get(ctx.Request.Context(), join.Secret)
	if !errors.Is(err, memcached.ErrNotFound) && err != nil {
		logger.Error("Error getting party from cache", zap.String("secret", join.Secret), zap.Any("player", player), zap.Error(err))

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{
			Message: "Internal server error. Please try again later.",
		})
		return
	}
	if errors.Is(err, memcached.ErrNotFound) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Error{
			Message: "The party does not exist",
		})
		return
	}

	var party response.Party

	// If the player is already in a party, publish a party
	// update event to all members and leave it
	currentItem, _ := cache.Get(ctx.Request.Context(), player.Hash())
	if currentItem != nil {
		// Close the channel if the player has subscribed to the party events
		channel, ok := subscriptions[player]
		if ok {
			close(channel)
		}

		party.Remove(player)
		party.Update(cache)

		json.Unmarshal(currentItem.Value, &party)
		flow.PublishAsync(party.JoinSecret, party)
	}

	// Use the same party object to retrieve the requested party
	json.Unmarshal(item.Value, &party)

	party.Add(player)

	bytes, _ := json.Marshal(party)
	err = cache.Set(ctx.Request.Context(), player.Hash(), bytes, 0, 0)
	if err != nil {
		logger.Error("Error setting party", zap.Any("player", player), zap.Error(err))

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{
			Message: "Internal server error. Please try again later.",
		})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusAccepted, party)
}
