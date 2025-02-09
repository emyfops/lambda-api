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
	Name: "lambda_rpc_logged_in_users",
	Help: "Total number of logged in users",
}, []string{"version"})

// JoinParty godoc
//
//	@Summary	Join a party
//	@Tags		Party
//	@Accept		json
//	@Produce	json
//	@Param		id	body		string	true	"Party ID"
//	@Success	202	{object}	response.Party
//	@Failure	400	{object}	response.ValidationError
//	@Failure	404	{object}	response.Error
//	@Router		/party/join [put]
//	@Security	ApiKeyAuth
func JoinParty(ctx *gin.Context, cache *memcache.Client) {
	var party response.Party
	var join request.JoinParty
	if err := ctx.Bind(&join); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ValidationError{
			Message: "Required fields are missing or invalid",
			Errors:  err.Error(),
		})
		return
	}

	if len(party.Players) >= party.Settings.MaxPlayers {
		ctx.AbortWithStatusJSON(http.StatusForbidden, response.Error{
			Message: "The party is full",
		})
		return
	}

	player := ctx.MustGet("player").(response.Player)

	_, err := cache.Get(player.Hash())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Error{
			Message: "The party does not exist",
		})
	}

	party.Add(player)

	bytes, _ := json.Marshal(party)
	cache.Set(&memcache.Item{Key: player.Hash(), Value: bytes})

	ctx.AbortWithStatusJSON(http.StatusAccepted, party)
	loggedInTotal.WithLabelValues("v1").Inc()
}
