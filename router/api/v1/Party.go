package v1

import (
	"github.com/Edouard127/lambda-rpc/internal/util"
	"github.com/Edouard127/lambda-rpc/pkg/cache"
	"github.com/Edouard127/lambda-rpc/pkg/party"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Persistent memory map to store the parties
// Owner -> Party
var partyMap = cache.NewPersistentMemoryCache[util.Player, *party.Party](0)

func CreateParty(ctx *gin.Context) {
	username := ctx.GetString("username")
	hash := ctx.GetString("hash")

	// Check if the user is allowed to create a party
	player := util.IsConnected(username, hash)
	if player == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "You either have an invalid account or the hash has expired, please reconnect to the server",
		})
	}

	// Check if the user is already in a party
	if _, exists := partyMap.Get(*player); exists {
		ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": "You are already in a party",
		})
	}

	// Create a party
	party := party.New("Party", *player) // TODO: Being able to name parties ? i don't think so
	partyMap.Set(*player, party)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Party created",
		"party":   party,
	})
}
