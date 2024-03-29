package endpoints

import (
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models"
	"github.com/Edouard127/lambda-rpc/pkg/io"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Persistent memory map to store the parties
// Party ID -> Party
var partyMap = io.NewPersistentMemoryCache[string, *models.Party](0)

// CreateParty godoc
// @BasePath /api/v1
// @Summary Create a new party
// @Description Create a new party
// @Tags Party
// @Accept json
// @Produce json
// @Success 201 {object} models.Party
// @Router /party/create [post]
// @Security Bearer
func CreateParty(ctx *gin.Context) {
	username := ctx.GetString("username")
	hash := ctx.GetString("hash")
	token := ctx.GetString("token")

	// Check if the user is allowed to create a party
	player := models.GetPlayer(username, hash, token)
	if player == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "You either have an invalid account or the hash has expired, please reconnect to the server",
		})
	}

	// Check if the user is already in a party
	if _, exists := partyMap.Get(hash); exists {
		ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": "You are already in a party",
		})
	}

	// Create a party
	p := models.New(*player)
	partyMap.Set(p.ID, p)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Party created",
		"party":   p,
	})
}
