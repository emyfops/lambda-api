package middlewares

import (
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models"
	"github.com/Edouard127/lambda-rpc/pkg/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func CheckAuth(ctx *gin.Context) {
	authorization := ctx.GetHeader("Authorization")
	token := strings.Split(authorization, "Bearer ")
	if len(token) != 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "You either have an invalid account or the hash has expired, please reconnect to the server",
		})
		return
	}

	jwt, err := auth.ParseJwtToken(token[1])
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "You either have an invalid account or the hash has expired, please reconnect to the server",
		})
		return
	}

	var player models.Player
	err = auth.ParseToStruct(jwt, &player)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "You either have an invalid account or the hash has expired, please reconnect to the server",
		})
		return
	}

	ctx.Set("player", player)
	ctx.Next()
}
