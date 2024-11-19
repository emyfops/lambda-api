package endpoints

import (
	"context"
	"github.com/Edouard127/lambda-api/internal/app/gonic"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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
func DeleteParty(ctx *gin.Context, client *redis.Client) {
	var party response.Party
	player := gonic.MustGet[response.Player](ctx, "player")

	err := client.HGetAll(context.Background(), player.String()).Scan(&party)
	if err != nil { // todo: change this
		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Error{
			Message: "You are not in a party",
		})
		return
	}

	if party.Leader != player {
		ctx.AbortWithStatusJSON(http.StatusForbidden, response.Error{
			Message: "You are not the leader of the party",
		})
		return
	}

	client.Del(context.Background(), player.String())

	ctx.AbortWithStatus(http.StatusNoContent)
	partyCountTotal.WithLabelValues("v1").Dec()
	loggedInTotal.WithLabelValues("v1").Sub(float64(len(party.Players)))
}
