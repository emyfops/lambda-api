package middlewares

import (
	"github.com/Edouard127/lambda-api/internal/app/auth"
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

	jwt, err := auth.ParseString(token[1][:len(token[1])-1]) // TODO: Temporary fix for the trailing quotation mark
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "The JWT token is invalid, please provide a valid JWT token",
		})
		return
	}

	var player response.Player
	err = auth.ParseStruct(jwt, &player)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "You either have an invalid account or the hash has expired, please reconnect to the server",
		})
		return
	}

	ctx.Set("player", player)
	ctx.Next()
}
