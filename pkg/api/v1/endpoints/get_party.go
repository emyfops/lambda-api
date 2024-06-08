package endpoints

import (
	"github.com/Edouard127/lambda-rpc/internal/app/auth"
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetParty godoc
// @BasePath /api/v1
// @Summary Get the party of the player
// @Tags Party
// @Accept json
// @Produce json
// @Success 200 {object} response.Party
// @Router /party [get]
// @Security Bearer
func GetParty(ctx *gin.Context) {
	player := auth.GinMustGet[response.Player](ctx, "player")

	id, exists := playerMap.Get(player)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Player not found",
		})
	}

	party, exists := partyMap.Get(id)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Party not found",
		})
	}

	ctx.AbortWithStatusJSON(http.StatusOK, party)
}
