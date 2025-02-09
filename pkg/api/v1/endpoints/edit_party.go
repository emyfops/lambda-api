package endpoints

import (
	"encoding/json"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/request"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
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
func EditParty(ctx *gin.Context, cache *memcache.Client) {
	var party response.Party
	var settings request.Settings
	if err := ctx.Bind(&settings); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ValidationError{
			Message: "Required fields are missing or invalid",
			Errors:  err.Error(),
		})
		return
	}

	player := ctx.MustGet("player").(response.Player)

	item, err := cache.Get(player.Hash())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Error{
			Message: "You are not in a party",
		})
		return
	}

	json.Unmarshal(item.Value, &party)

	if party.Leader != player {
		ctx.AbortWithStatusJSON(http.StatusForbidden, response.Error{
			Message: "You are not the leader of the party",
		})
		return
	}

	party.Settings = settings

	bytes, _ := json.Marshal(party)
	cache.Set(&memcache.Item{Key: player.Hash(), Value: bytes})

	ctx.AbortWithStatusJSON(http.StatusAccepted, party)
}
