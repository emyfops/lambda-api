package endpoints

import "github.com/gin-gonic/gin"

// ListParties godoc
// @BasePath /api/v1
// @Summary List all parties
// @Description List all parties
// @Tags Party
// @Accept json
// @Produce json
// @Success 200 {array} response.Party
// @Router /party/list [get]
func ListParties(ctx *gin.Context) {
	// TODO
	ctx.AbortWithStatusJSON(501, gin.H{
		"message": "Not implemented",
	})
	return
}
