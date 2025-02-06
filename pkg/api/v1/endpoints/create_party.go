package endpoints

import (
	"context"
	"errors"
	"github.com/Edouard127/lambda-api/internal/app/gonic"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/request"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/redis/go-redis/v9"
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
func CreateParty(ctx *gin.Context, client *redis.Client) {
	var settings request.Settings
	if err := ctx.Bind(&settings); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ValidationError{
			Message: "Required fields are missing or invalid",
			Errors:  err.Error(),
		})
		return
	}

	player := gonic.MustGet[response.Player](ctx, "player")

	_, err := client.Get(context.Background(), player.String()).Result()
	if !errors.Is(err, redis.Nil) {
		// We should only check against redis.Nil
		// If we get another error it most likely means
		// that something went wrong, either the party
		// doesn't exist and redis is acting up, or the
		// party was lost and needs to be recreated
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
	client.HSet(context.Background(), player.String(), party)

	ctx.AbortWithStatusJSON(http.StatusCreated, party)
	partyCountTotal.WithLabelValues("v1").Inc()
	loggedInTotal.WithLabelValues("v1").Inc()
}
