package endpoints

import (
	"encoding/json"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetParty 	godoc
//
//	@Summary	Get the party of the player
//	@Tags		Party
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	response.Party
//	@Failure	404	{object}	response.Error
//	@Router		/party [get]
//	@Security	ApiKeyAuth
func GetParty(ctx *gin.Context, cache *memcache.Client) {
	var party response.Party
	player := ctx.MustGet("player").(response.Player)

	item, err := cache.Get(player.Hash())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Error{
			Message: "You are not in a party",
		})
		return
	}

	json.Unmarshal(item.Value, &party)

	ctx.AbortWithStatusJSON(http.StatusOK, party)
}
