package endpoints

import (
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models"
	"github.com/Edouard127/lambda-rpc/pkg/auth"
	"github.com/Edouard127/lambda-rpc/pkg/io"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Persistent memory map to store the parties
// Party ID -> Party
var partyMap = io.NewPersistentMemoryCache[string, *models.Party](0)

// Player -> Party ID
var playerMap = io.NewPersistentMemoryCache[models.Player, string](0)

// CreateParty godoc
// @BasePath /api/v1
// @Summary Create a new party
// @Description Create a new party
// @Tags Party
// @Accept json
// @Produce json
// @Param Settings body models.Settings false "Settings"
// @Success 201 {object} models.Party
// @Failure 409 {object} models.Party
// @Router /party/create [post]
// @Security Bearer
func CreateParty(ctx *gin.Context) {
	var settings *models.Settings
	err := ctx.BindJSON(&settings)
	if err != nil {
		return
	}

	player := auth.GinMustGet[models.Player](ctx, "player")

	// Check if the player is already in a party
	if partyID, exists := playerMap.Get(player); exists {
		party, _ := partyMap.Get(partyID)
		ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": "You are already in a party",
			"party":   party,
		})
	}

	party := models.NewWithSettings(player, settings)

	partyMap.Set(party.ID, party)
	playerMap.Set(player, party.ID) // Reverse mapping

	ctx.AbortWithStatusJSON(http.StatusCreated, party)
	return
}
