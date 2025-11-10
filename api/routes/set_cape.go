package routes

import (
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Edouard127/lambda-api/api/models"
	"github.com/Edouard127/lambda-api/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

var capeList = make(map[string]struct{})

func init() {
	r, err := http.Get("https://raw.githubusercontent.com/emyfops/lambda-assets/refs/heads/master/capes.txt")
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

func SetCape(ctx *fiber.Ctx) error {
	logger := internal.MustGetState[*slog.Logger]("logger")
	cache := internal.MustGetState[*redis.Client]("cache")
	player := ctx.Locals("player").(models.Player)

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
