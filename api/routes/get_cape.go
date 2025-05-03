package routes

import (
	"errors"
	"github.com/Edouard127/lambda-api/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"net/http"
)

func GetCape(ctx *fiber.Ctx) error {
	logger := internal.MustGetState[*slog.Logger]("logger")
	cache := internal.MustGetState[*redis.Client]("cache")

	id := ctx.Query("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid player id")
	}

	cape, err := cache.Get(ctx.UserContext(), id).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fiber.NewError(http.StatusNotFound, "this player does not have a cape")
		}

		logger.Error("Error getting player query from cache", slog.String("id", id), slog.Any("error", err))

		return fiber.NewError(http.StatusInternalServerError, "internal server error")
	}

	return ctx.JSON(fiber.Map{"id": uid, "type": cape})
}
