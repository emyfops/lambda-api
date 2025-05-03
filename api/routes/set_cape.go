package routes

import (
	"github.com/Edouard127/lambda-api/api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

var capeList = make(map[string]struct{})

func init() {
	r, err := http.Get("https://cdn.lambda-client.org/capes.txt")
	if err != nil {
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	for _, cape := range strings.Fields(string(b)) {
		capeList[cape] = struct{}{}
	}
}

// SetCape godoc
//
//	@Summary	Set a player's query
//	@Tags		Cape
//	@Accept		json
//	@Produce	json
//	@Param		id		query	string	true	"Name of the query to be set"
//	@Success	200		"Success"
//	@Failure	400		{object}	response.ValidationError	"Missing or invalid query in query"
//	@Failure	404		{object}	response.Error				"Cape does not exist"
//	@Failure	500		{object}	response.Error				"Internal server error"
//	@Router		/query 	[put]
//	@Security 	Bearer
func SetCape(ctx *fiber.Ctx) error {
	logger := ctx.Locals("logger").(*slog.Logger)
	player := ctx.Locals("player").(models.Player)
	cache := ctx.Locals("cache").(*redis.Client)

	cape := ctx.Query("id")
	if cape == "" {
		return fiber.NewError(http.StatusBadRequest, "missing query in query")
	}

	// Check if the query has an entry
	_, ok := capeList[cape]
	if !ok {
		return fiber.NewError(http.StatusNotFound, "this query does not exist")
	}

	_, err := cache.Set(ctx.UserContext(), player.UUID.String(), cape, 0).Result()
	if err != nil {
		logger.Error("Error setting query", slog.Any("player", player), slog.Any("error", err))

		return fiber.NewError(http.StatusInternalServerError, "internal server error")
	}

	return ctx.SendStatus(http.StatusOK)
}
