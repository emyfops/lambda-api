package endpoints

import (
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models/request"
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models/response"
	"github.com/Edouard127/lambda-rpc/pkg/auth"
	"github.com/Edouard127/lambda-rpc/pkg/io"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Persistent memory map to store the parties
// Party ID -> Party
var partyMap = io.NewPersistentMemoryCache[string, *response.Party](0)

// Player -> Party ID
var playerMap = io.NewPersistentMemoryCache[response.Player, string](0)

// CreateParty godoc
// @BasePath /api/v1
// @Summary Create a new party
// @Description Create a new party
// @Tags Party
// @Accept json
// @Produce json
// @Param Settings body request.Settings false "Settings"
// @Success 201 {object} response.Party
// @Router /party/create [post]
// @Security Bearer
func CreateParty(ctx *gin.Context) {
	var settings request.Settings

	if err := ctx.ShouldBind(&settings); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	player := auth.GinMustGet[response.Player](ctx, "player")

	// Check if the player is already in a party
	if partyID, exists := playerMap.Get(player); exists {
		party, _ := partyMap.Get(partyID)
		ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": "You are already in a party",
			"party":   party,
		})
	}

	party := response.NewWithSettings(player, &settings)

	partyMap.Set(party.ID, party)
	playerMap.Set(player, party.ID) // Reverse mapping

	ctx.AbortWithStatusJSON(http.StatusCreated, party)
	return
}
