package routes

import (
	"github.com/Edouard127/lambda-api/api/models"
	"github.com/Edouard127/lambda-api/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"regexp"
	"time"
)

// Login allows a player to log in to the server using a Minecraft username and Mojang session hash
//
// Refer to the project README for more information on that process
func Login(ctx *fiber.Ctx) error {
	var login models.Authentication

	err := ctx.BodyParser(&login)
	if err != nil {
		return fiber.NewError(http.StatusUnprocessableEntity, "required fields are missing or invalid")
	}

	ok, err := regexp.MatchString("^[a-zA-Z0-9_]{2,16}$", login.Username)
	if !ok || err != nil {
		return fiber.NewError(http.StatusUnprocessableEntity, "required fields are missing or invalid")
	}

	player, err := models.GetPlayer(login.Username, login.Hash)
	if err != nil {
		return fiber.NewError(http.StatusUnauthorized, "invalid credentials")
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
