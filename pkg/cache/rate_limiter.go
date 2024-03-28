package cache

import (
	"github.com/gin-gonic/gin"
	"time"
)

type RateLimiter struct {
	Time   time.Duration
	afterN int
	cache  *MemoryCache[string, int]
}

func NewRateLimiter(time time.Duration, threshold int) *RateLimiter {
	return &RateLimiter{
		Time:   time,
		afterN: threshold,
		cache:  NewTempMemoryCache[string, int](time, time, 0),
	}
}

func (rt *RateLimiter) Handle(ctx *gin.Context) {
	ip := ctx.ClientIP()
	if val, ok := rt.cache.Get(ip); ok {
		if val >= rt.afterN {
			ctx.JSON(429, gin.H{
				"message": "Rate limit exceeded",
			})
			return
		}
		rt.cache.Set(ip, val+1)
	} else {
		rt.cache.Set(ip, 1)
	}
	ctx.Next()
}
