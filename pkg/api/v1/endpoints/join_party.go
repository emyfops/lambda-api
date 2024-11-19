package endpoints

import (
	"context"
	"github.com/Edouard127/lambda-api/internal/app/gonic"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/request"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"net/http"
)

var (
	loggedInTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "lambda_rpc_logged_in_users",
		Help: "Total number of logged in users",
	}, []string{"version"})
)

func init() {
	prometheus.MustRegister(loggedInTotal)
}

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
func JoinParty(ctx *gin.Context, client *redis.Client) {
	var party response.Party
	var join request.JoinParty
	if err := ctx.Bind(&join); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ValidationError{
			Message: "Required fields are missing or invalid",
			Errors:  err.Error(),
		})
		return
	}

	player := gonic.MustGet[response.Player](ctx, "player")

	err := client.HGetAll(context.Background(), player.String()).Scan(&party)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Error{
			Message: "The party does not exist",
		})
	}

	party.Add(player)
	client.HSet(context.Background(), player.String(), "players", party.Players)

	ctx.AbortWithStatusJSON(http.StatusAccepted, party)
	loggedInTotal.WithLabelValues("v1").Inc()
}
