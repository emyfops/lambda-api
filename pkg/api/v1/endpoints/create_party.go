package endpoints

import (
	"encoding/json"
	"errors"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/yeqown/memcached"
	"go.uber.org/zap"
	"net/http"
)

var (
	partyCountTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "lambda_api_party_count_total",
		Help: "Total number of parties",
	}, []string{"version"})
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

	partyCountTotal.WithLabelValues("v1").Inc()
	loggedInTotal.WithLabelValues("v1").Inc()

	ctx.AbortWithStatusJSON(http.StatusCreated, party)
	return

err500:
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{
		Message: "Internal server error. Please try again later.",
	})
	return
}
