package endpoints

import (
	"encoding/json"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	"github.com/bradfitz/gomemcache/memcache"
	eventbus "github.com/dtomasi/go-event-bus/v3"
	"github.com/gin-gonic/gin"
	"net/http"
)

var flow = eventbus.NewEventBus()

func PartyListen(ctx *gin.Context, cache *memcache.Client) {
	player := ctx.MustGet("player").(response.Player)

	item, err := cache.Get(player.Hash())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Error{
			Message: "You are not in a party",
		})
		return
	}

	var party response.Party
	json.Unmarshal(item.Value, &party)

	// Setup mandatory headers for SSE
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")

	partyChannel := flow.Subscribe("1")

	for {
		event := <-partyChannel

		//newParty := event.Data.(response.Party)
		ctx.SSEvent("1", event)

		// Flush the data immediately instead of buffering it for later.
		ctx.Writer.Flush()
	}
}
