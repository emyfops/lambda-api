package middlewares

import (
	"encoding/json"
	"github.com/Edouard127/lambda-api/api/models"
	"github.com/Edouard127/lambda-api/internal"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

var MinecraftCheck = jwtware.New(jwtware.Config{
	ErrorHandler: ErrorHandler,
	SuccessHandler: func(ctx *fiber.Ctx) error {
		token := ctx.Locals("user").(*jwt.Token)

		parsed := token.Claims.(jwt.MapClaims)["data"]
		bytes, _ := json.Marshal(parsed)

		var player models.Player
		json.Unmarshal(bytes, &player)

		ctx.Locals("player", player)
		return ctx.Next()
	},
	SigningKey: jwtware.SigningKey{
		JWTAlg: "RS256",
		Key:    internal.PrivateKey.Public(),
	},
})

// DiscordCheck sits between the authorization middleware and the handler function
// and ensure that the player has linked their discord account
func DiscordCheck(ctx *fiber.Ctx) error {
	player := ctx.Locals("player").(models.Player)

	if !player.HasDiscord() {
		return fiber.NewError(http.StatusUnauthorized, "you did not link your discord account")
	}

	return ctx.Next()
}
