package routes

import (
	"github.com/Edouard127/lambda-api/api/models"
	"github.com/Edouard127/lambda-api/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

// LinkDiscord links a discord account to an existing bearer token
func LinkDiscord(ctx *fiber.Ctx) error {
	player := ctx.Locals("player").(models.Player)

	var link models.DiscordLink

	err := ctx.BodyParser(&link)
	if err != nil {
		return fiber.NewError(http.StatusUnprocessableEntity, "required fields are missing or invalid")
	}

	err = models.GetDiscord(link.Token, &player)
	if err != nil {
		return fiber.NewError(http.StatusUnauthorized, "invalid discord token")
	}

	claims := jwt.MapClaims{
		"nbf":  time.Now().Unix(),
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"data": player,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signed, err := token.SignedString(internal.PrivateKey)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "failed to create token")
	}

	return ctx.JSON(fiber.Map{
		"access_token": signed,
		"expires_in":   int64(time.Hour * 24),
		"token_type":   "Bearer",
	})
}
