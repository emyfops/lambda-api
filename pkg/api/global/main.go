package global

import (
	"github.com/Edouard127/lambda-api/internal/app/healthcheck"
	"github.com/alexliesenfeld/health"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"
)

func Register(client *redis.Client, router *gin.Engine) {
	global := router.Group("/api")

	// Initialize the healthcheck handler
	checker := health.NewChecker(
		// 2 second TTL
		health.WithCacheDuration(2*time.Second),

		// Global timeout of 5 seconds
		health.WithTimeout(5*time.Second),

		// Mojang API readiness check
		health.WithPeriodicCheck(
			60*time.Second,
			time.Second,
			health.Check{
				Name:  "http-connection-mojang-session",
				Check: healthcheck.HTTPGetCheck("https://sessionserver.mojang.com/session/minecraft/hasJoined"),
			},
		),

		// Redis readiness check
		health.WithPeriodicCheck(
			5*time.Second,
			time.Second,
			health.Check{
				Name:  "redis-connection",
				Check: healthcheck.RedisCheck(client),
			},
		),
	)

	global.GET("/health", gin.WrapF(health.NewHandler(checker)))
}
