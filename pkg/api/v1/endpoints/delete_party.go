package endpoints

import (
	"encoding/json"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
	"net/http"
)

// DeleteParty 	godoc
//
//	@Summary	Delete an existing party
//	@Tags		Party
//	@Accept		json
//	@Produce	json
//	@Success	204
//	@Failure	403	{object}	response.Error
//	@Failure	404	{object}	response.Error
//	@Router		/party/delete [delete]
//	@Security	ApiKeyAuth
func DeleteParty(ctx *gin.Context, cache *memcache.Client) {
	player := ctx.MustGet("player").(response.Player)

	item, _ := cache.Get(player.Hash())
	if item == nil {
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

	cache.Delete(player.Hash())
	cache.Delete(party.JoinSecret)

	partyCountTotal.WithLabelValues("v1").Dec()
	loggedInTotal.WithLabelValues("v1").Sub(float64(len(party.Players)))

	ctx.AbortWithStatus(http.StatusNoContent)
}
