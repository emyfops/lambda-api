package endpoints

import (
	"encoding/json"
	"errors"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/request"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
)

var (
	partyCountTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "lambda_rpc_party_count_total",
		Help: "Total number of parties",
	}, []string{"version"})
)

// CreateParty 	godoc
//
//	@Summary	Create a new party
//	@Tags		Party
//	@Accept		json
//	@Produce	json
//	@Param		settings	body		request.Settings	false	"Party configuration"
//	@Success	201			{object}	response.Party
//	@Failure	400			{object}	response.ValidationError
//	@Failure	409			{object}	response.Error
//	@Router		/party/create [post]
//	@Security	ApiKeyAuth
func CreateParty(ctx *gin.Context, cache *memcache.Client) {
	var settings request.Settings
	if err := ctx.Bind(&settings); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ValidationError{
			Message: "Required fields are missing or invalid",
			Errors:  err.Error(),
		})
		return
	}

	player := ctx.MustGet("player").(response.Player)

	_, err := cache.Get(player.Hash())
	if !errors.Is(err, memcache.ErrCacheMiss) && err != nil {
		// We should only check against memcache.ErrCacheMiss
		// If we get another error it most likely means that something went wrong, either the
		// cache is acting up, or the party was lost and needs to be recreated
		ctx.AbortWithStatusJSON(http.StatusConflict, response.Error{
			Message: "You are already in a party",
		})
		return
	}

	party := response.NewParty(player, &settings)

	// We are using redis for the regular mapping and for
	// the ability of scaling horizontally without losing
	// data if containers are scaled down
	//
	// Mapping: Party ID -> Party struct
	bytes, _ := json.Marshal(party)
	cache.Set(&memcache.Item{Key: player.Hash(), Value: bytes})

	ctx.AbortWithStatusJSON(http.StatusCreated, party)
	partyCountTotal.WithLabelValues("v1").Inc()
	loggedInTotal.WithLabelValues("v1").Inc()
}
