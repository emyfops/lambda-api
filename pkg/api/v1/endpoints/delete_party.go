package endpoints

import (
	"github.com/Edouard127/lambda-rpc/internal/app/auth"
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// DeleteParty godoc
// @BasePath /api/v1
// @Summary Delete an existing party
// @Tags Party
// @Accept json
// @Produce json
// @Success 204
// @Router /party/delete [delete]
// @Security Bearer
func DeleteParty(ctx *gin.Context) {
	player := auth.GinMustGet[response.Player](ctx, "player")

	partyID, exists := playerMap.Get(player)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "You are not in a party",
		})
	}

	party, exists := partyMap.Get(partyID)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Party not found",
		})
	}

	if party.Leader != player {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "You are not the leader of the party",
		})
		return
	}

	partyMap.Delete(partyID)
	playerMap.Delete(player)

	ctx.AbortWithStatus(http.StatusNoContent)
}
