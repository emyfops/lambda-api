package middlewares

import (
	"github.com/Edouard127/lambda-rpc/internal/app/io"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var cache = io.NewTempMemoryCache[string, int](time.Second*10, time.Second*5, 0)

func RateLimit(ctx *gin.Context) {
	ip := ctx.ClientIP()
	if n, ok := cache.Get(ip); ok && n > 5 {
		ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"message": "You are being rate limited, please try again later",
		})
		return
	} else {
		cache.Set(ip, n+1)
		ctx.Next()
	}
}
