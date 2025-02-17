package middlewares

import (
	"github.com/Edouard127/lambda-api/internal"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func CheckAuth(ctx *gin.Context) {
	authorization := ctx.GetHeader("Authorization")
	token := strings.Split(authorization, "Bearer ")
	if len(token) != 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "The Authorization header is missing or invalid, please provide a valid JWT token with the Bearer prefix",
		})
		return
	}

	signed, err := internal.ParseJwt(token[1])
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "The JWT token is invalid, please provide a valid JWT token",
		})
		return
	}

	var player response.Player
	err = internal.ParseStructJwt(signed, &player)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred while parsing the JWT token",
		})
		return
	}

	ctx.Set("player", player)
	ctx.Next()
}
