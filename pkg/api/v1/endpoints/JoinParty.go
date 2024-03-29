package endpoints

import "github.com/gin-gonic/gin"

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

}
