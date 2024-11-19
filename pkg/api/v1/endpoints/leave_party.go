package endpoints

import (
	"context"
	"github.com/Edouard127/lambda-api/internal/app/gonic"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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
func LeaveParty(ctx *gin.Context, client *redis.Client) {
	var party response.Party
	player := gonic.MustGet[response.Player](ctx, "player")

	err := client.HGetAll(context.Background(), player.String()).Scan(&party)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Error{
			Message: "You are not in a party",
		})
		return
	}

	client.Del(context.Background(), player.String())

	ctx.AbortWithStatus(http.StatusAccepted)
	loggedInTotal.WithLabelValues("v1").Dec()
}
