package global

import (
	"github.com/Edouard127/lambda-api/internal"
	"github.com/alexliesenfeld/health"
	"github.com/gin-gonic/gin"
	"time"
)

func Register(router *gin.Engine) {
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
				Check: internal.HTTPGetCheck("https://sessionserver.mojang.com/session/minecraft/hasJoined"),
			},
		),
	)

	global.GET("/health", gin.WrapF(health.NewHandler(checker)))
}
