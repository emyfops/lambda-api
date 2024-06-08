package endpoints

import (
	"github.com/Edouard127/lambda-rpc/internal/app/auth"
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models/request"
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// EditParty godoc
// @BasePath /api/v1
// @Summary Edit a party
// @Tags Party
// @Accept json
// @Produce json
// @Param Settings body request.Settings false "Settings"
// @Success 202 {object} response.Party
// @Failure 400 {object} response.Error
// @Failure 403 {object} response.ValidationError
// @Failure 404 {object} response.Error
// @Router /party/edit [patch]
// @Security Bearer
func EditParty(ctx *gin.Context) {
	var settings request.Settings
	if err := ctx.Bind(&settings); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ValidationError{
			Message: "Required fields are missing or invalid",
			Errors:  err.Error(),
		})
		return
	}

	player := auth.GinMustGet[response.Player](ctx, "player")

	partyID, exists := playerMap.Get(player)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Error{
			Message: "You are not in a party",
		})
	}

	party, exists := partyMap.Get(partyID)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Error{
			Message: "The party does not exist",
		})
	}

	if party.Leader != player {
		ctx.AbortWithStatusJSON(http.StatusForbidden, response.Error{
			Message: "You are not the leader of the party",
		})
		return
	}

	party.Settings = settings

	ctx.AbortWithStatusJSON(http.StatusAccepted, party)
}
