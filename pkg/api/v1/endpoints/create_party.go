package endpoints

import (
	"github.com/Edouard127/lambda-rpc/internal/app/auth"
	"github.com/Edouard127/lambda-rpc/internal/app/io"
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models/request"
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Player -> &Party ID
var playerMap = io.NewCache[response.Player, string]()

// Reverse mapping of playerMap
// &Party ID -> &Party
var partyMap = io.NewCache[string, *response.Party]()

// CreateParty godoc
// @BasePath /api/v1
// @Summary Create a new party
// @Tags Party
// @Accept json
// @Produce json
// @Param Settings body request.Settings false "Settings"
// @Success 201 {object} response.Party
// @Failure 400 {object} response.ValidationError
// @Failure 409 {object} response.Error
// @Router /party/create [post]
// @Security Bearer
func CreateParty(ctx *gin.Context) {
	var settings request.Settings
	if err := ctx.Bind(&settings); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ValidationError{
			Message: "Required fields are missing or invalid",
			Errors:  err.Error(),
		})
		return
	}

	player := auth.GinMustGet[response.Player](ctx, "player")

	_, exists := playerMap.Get(player)
	if exists {
		ctx.AbortWithStatusJSON(http.StatusConflict, response.Error{
			Message: "You are already in a party",
		})
	}

	party := response.NewWithSettings(player, &settings)

	partyMap.Set(party.ID, party, -1)
	playerMap.Set(player, party.ID, -1)

	ctx.AbortWithStatusJSON(http.StatusCreated, party)
}
