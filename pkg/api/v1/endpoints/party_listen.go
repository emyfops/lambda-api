package endpoints

import (
	"encoding/json"
	"errors"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	eventbus "github.com/dtomasi/go-event-bus/v3"
	"github.com/gin-gonic/gin"
	"github.com/yeqown/memcached"
	"go.uber.org/zap"
	"net/http"
)

var subscriptions = make(map[response.Player]eventbus.EventChannel)
var flow = eventbus.NewEventBus()

// PartyListen godoc
//
//	@Summary	Listen for party updates via SSE
//	@Tags		Party
//	@Accept		json
//	@Produce	text/event-stream
//	@Success	200	"Streaming party events"
//	@Failure	404	{object}	response.Error	"You are not in a party"
//	@Failure	500	{object}	response.Error	"Internal server error"
//	@Router		/party/listen [get]
//	@Security 	Bearer
func PartyListen(ctx *gin.Context, cache memcached.Client) {
	logger := ctx.MustGet("logger").(*zap.Logger)
	player := ctx.MustGet("player").(response.Player)

	item, err := cache.Get(ctx.Request.Context(), player.Hash())
	if !errors.Is(err, memcached.ErrNotFound) && err != nil {
		logger.Error("Error getting party from cache", zap.Any("player", player), zap.Error(err))

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{
			Message: "Internal server error. Please try again later.",
		})
		return
	}
	if errors.Is(err, memcached.ErrNotFound) {
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

	partyChannel := flow.Subscribe(party.JoinSecret)
	subscriptions[player] = partyChannel

	for {
		event, ok := <-partyChannel
		if !ok {
			return
		}

		ctx.SSEvent("data", event.Data)
		ctx.Writer.Flush()
	}
}
