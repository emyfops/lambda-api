package middlewares

import (
	"github.com/Edouard127/lambda-rpc/internal/app/memory"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"time"
)

var cache = memory.NewCache[string, *RateLimiterItem]()

type Options struct {
	// The maximum number of requests allowed in the duration.
	Limit rate.Limit

	// The duration in which the maximum number of requests is allowed.
	Duration time.Duration

	// Burst is the maximum number of requests that can be made in a short amount of time.
	Burst int
}

type RateLimiterItem struct {
	Limiter       *rate.Limiter
	LastSeenAt    time.Time
	IllegalAccess int
}

func RateLimiter(options Options) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		instance, exists := cache.Get(ip)
		if !exists ||
			time.Since((*instance).LastSeenAt).Milliseconds() >
				options.Duration.Milliseconds() {
			newItem(ip, &options)

			// It should be safe to reassign the pointer
			// because other owners have the pointer to the pointer
			instance, exists = cache.Get(ip)
		}

		if !(*instance).Limiter.Allow() ||
			time.Duration((*instance).Limiter.Reserve().Delay()) > 0 {
			ctx.AbortWithStatusJSON(429, gin.H{
				"message": "Too many requests",
			})
			return
		}

		(*instance).LastSeenAt = time.Now()

		ctx.Next()
	}
}

func newItem(ip string, options *Options) *RateLimiterItem {
	if _, exists := cache.Get(ip); exists {
		cache.Delete(ip)
	}

	instance := &RateLimiterItem{
		Limiter: rate.NewLimiter(options.Limit, options.Burst),
	}

	cache.Set(ip, instance, -1)

	return instance
}
