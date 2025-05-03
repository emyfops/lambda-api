package routes

import (
	"errors"
	"github.com/Edouard127/lambda-api/api/models"
	"github.com/Edouard127/lambda-api/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"net/http"
)

// GetCapes returns a list of player
func GetCapes(ctx *fiber.Ctx) error {
	logger := internal.MustGetState[*slog.Logger]("logger")
	cache := internal.MustGetState[*redis.Client]("cache")

	var lookup models.CapeLookup

	err := ctx.BodyParser(&lookup)
	if err != nil {
		return fiber.NewError(http.StatusUnprocessableEntity, "required fields are missing or invalid")
	}

	if len(lookup.Players) == 0 {
		return ctx.SendString("[]")
	}

	var ids = make([]string, 0)
	for _, id := range lookup.Players {
		ids = append(ids, id.String())
	}

	var capes []struct {
		Uuid uuid.UUID
		Type string
	}

	keys, err := cache.MGet(ctx.UserContext(), ids...).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.Error("Error getting player query from cache", slog.Any("error", err))

		return fiber.NewError(http.StatusInternalServerError, "internal server error")
	}

	for _, item := range keys {
		id, err := uuid.Parse(item.(string))
		if err != nil {
			continue
		}

		capes = append(capes, struct {
			Uuid uuid.UUID
			Type string
		}{
			Uuid: id,
			Type: item.(string),
		})
	}

	return ctx.JSON(capes)
}
