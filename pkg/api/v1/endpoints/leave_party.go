package endpoints

import (
	"github.com/Edouard127/lambda-rpc/internal/app/auth"
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// LeaveParty godoc
// @BasePath /api/v1
// @Summary Leave a party
// @Tags Party
// @Accept json
// @Produce json
// @Success 202
// @Failure 404 {object} response.Error
// @Router /party/leave [put]
// @Security Bearer
func LeaveParty(ctx *gin.Context) {
	player := auth.GinMustGet[response.Player](ctx, "player")

	id, exists := playerMap.Get(player)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Error{
			Message: "You are not in a party",
		})
		return
	}

	partyMap.Delete(*id)
	playerMap.Delete(player)
	loggedInTotal.WithLabelValues("v1").Dec()

	ctx.AbortWithStatus(http.StatusAccepted)
}
