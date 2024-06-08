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
// @Router /party/join [put]
// @Security Bearer
func JoinParty(ctx *gin.Context) {
	var join request.JoinParty
	if err := ctx.Bind(&join); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Validation error",
			"errors":  err.Error(),
		})
		return
	}

	player := auth.GinMustGet[response.Player](ctx, "player")

	party, exists := partyMap.Get(join.ID)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Party not found",
		})
	}

	party.Add(player)

	ctx.AbortWithStatusJSON(http.StatusAccepted, party)
	return
}
