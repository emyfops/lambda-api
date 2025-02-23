package endpoints

import (
	"encoding/json"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/request"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
//	@Param		id	body		request.JoinParty	true	"Party ID"
//	@Success	202	{object}	response.Party
//	@Failure	400	{object}	response.ValidationError
//	@Failure	404	{object}	response.Error
//	@Router		/party/join [put]
//	@Security	ApiKeyAuth
func JoinParty(ctx *gin.Context, cache *memcache.Client) {
	var join request.JoinParty

	err := ctx.Bind(&join)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ValidationError{
			Message: "Required fields are missing or invalid",
			Errors:  err.Error(),
		})
		return
	}

	player := ctx.MustGet("player").(response.Player)

	item, err := cache.Get(join.Secret)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Error{
			Message: "The party does not exist",
		})
		return
	}

	var party response.Party

	// If the player is already in a party, publish a party
	// update event to all members and leave it
	currentItem, _ := cache.Get(player.Hash())
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
	cache.Set(&memcache.Item{Key: player.Hash(), Value: bytes})

	ctx.AbortWithStatusJSON(http.StatusAccepted, party)
}
