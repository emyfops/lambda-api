package endpoints

import (
	"encoding/json"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
	"net/http"
)

// LeaveParty godoc
//
//	@Summary	Leave a party
//	@Tags		Party
//	@Accept		json
//	@Produce	json
//	@Success	202
//	@Failure	404	{object}	response.Error
//	@Router		/party/leave [put]
//	@Security	ApiKeyAuth
func LeaveParty(ctx *gin.Context, cache *memcache.Client) {
	player := ctx.MustGet("player").(response.Player)

	item, err := cache.Get(player.Hash())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Error{
			Message: "You are not in a party",
		})
		return
	}

	var party response.Party
	json.Unmarshal(item.Value, &party)

	party.Remove(player)
	flow.PublishAsync(party.JoinSecret, party)

	cache.Delete(player.Hash())

	loggedInTotal.WithLabelValues("v1").Dec()
	ctx.AbortWithStatus(http.StatusAccepted)
}
