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
// @Failure 404 {object} response.Error
// @Router /party [get]
// @Security Bearer
func GetParty(ctx *gin.Context) {
	player := auth.GinMustGet[response.Player](ctx, "player")

	id, exists := playerMap.Get(player)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Error{
			Message: "You are not in a party",
		})
	}

	party, exists := partyMap.Get(*id)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Error{
			Message: "The party does not exist",
		})
	}

	ctx.AbortWithStatusJSON(http.StatusOK, party)
}
