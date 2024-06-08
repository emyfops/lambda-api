package endpoints

import (
	"github.com/Edouard127/lambda-rpc/internal/app/auth"
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models/request"
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// JoinParty godoc
// @BasePath /api/v1
// @Summary Join a party
// @Tags Party
// @Accept json
// @Produce json
// @Param ID body string true "Party ID"
// @Success 202 {object} response.Party
// @Failure 400 {object} response.ValidationError
// @Failure 404 {object} response.Error
// @Router /party/join [put]
// @Security Bearer
func JoinParty(ctx *gin.Context) {
	var join request.JoinParty
	if err := ctx.Bind(&join); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ValidationError{
			Message: "Required fields are missing or invalid",
			Errors:  err.Error(),
		})
		return
	}

	player := auth.GinMustGet[response.Player](ctx, "player")

	party, exists := partyMap.Get(join.ID)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Error{
			Message: "The party does not exist",
		})
	}

	(*party).Add(player)

	ctx.AbortWithStatusJSON(http.StatusAccepted, party)
}
