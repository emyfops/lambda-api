package endpoints

import (
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models"
	"github.com/Edouard127/lambda-rpc/pkg/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

// JoinParty godoc
// @BasePath /api/v1
// @Summary Join a party
// @Description Join a party
// @Tags Party
// @Accept json
// @Produce json
// @Success 201 {object} models.Party
// @Router /party/join [post]
// @Security Bearer
func JoinParty(ctx *gin.Context) {
	player := auth.GinMustGet[models.Player](ctx, "player")

	id, exists := playerMap.Get(player)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Party not found",
		})
	}

	party, exists := partyMap.Get(id)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Party not found",
		})
	}

	party.Add(player)

	ctx.JSON(http.StatusCreated, party)
}
