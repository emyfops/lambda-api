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
// @Param ID query string true "Party ID"
// @Success 202 {object} models.Party
// @Router /party/join [put]
// @Security Bearer
func JoinParty(ctx *gin.Context) {
	player := auth.GinMustGet[models.Player](ctx, "player")

	party, exists := partyMap.Get(ctx.Query("ID"))
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Party not found",
		})
	}

	party.Add(player)
	ctx.AbortWithStatusJSON(http.StatusAccepted, party)
}
