package middlewares

import (
	"github.com/Edouard127/lambda-rpc/internal/util"
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func CheckAuth(ctx *gin.Context) {
	auth := ctx.GetHeader("Authorization")
	token := strings.Split(auth, "Bearer ")
	if len(token) != 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "You either have an invalid account or the hash has expired, please reconnect to the server",
		})
		return
	}

	jwt, err := util.ParseJwtToken(token[1])
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "You either have an invalid account or the hash has expired, please reconnect to the server",
		})
		return
	}

	var player models.Player
	err = util.JwtToStruct(jwt, &player)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "You either have an invalid account or the hash has expired, please reconnect to the server",
		})
		return
	}

	ctx.Set("player", player)
	ctx.Next()
}
