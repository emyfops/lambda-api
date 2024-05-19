package endpoints

import (
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models/request"
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models/response"
	"github.com/Edouard127/lambda-rpc/pkg/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

// EditParty godoc
// @BasePath /api/v1
// @Summary Edit a party
// @Description Edit a party
// @Tags Party
// @Accept json
// @Produce json
// @Param Settings body request.Settings false "Settings"
// @Success 202 {object} response.Party
// @Router /party/edit [patch]
// @Security Bearer
func EditParty(ctx *gin.Context) {
	var settings request.Settings

	if err := ctx.Bind(&settings); err != nil {
		return
	}

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

	party.Settings = settings

	// TODO: Find a way to notify the party members that the settings have changed

	ctx.AbortWithStatusJSON(http.StatusAccepted, party)
	return
}
