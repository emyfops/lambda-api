package endpoints

import (
	"context"
	"github.com/Edouard127/lambda-api/internal/app/gonic"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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
func GetParty(ctx *gin.Context, client *redis.Client) {
	var party response.Party
	player := gonic.MustGet[response.Player](ctx, "player")

	err := client.HGetAll(context.Background(), player.String()).Scan(&party)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Error{
			Message: "You are not in a party",
		})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, party)
}
