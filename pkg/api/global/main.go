package global

import (
	"github.com/Edouard127/lambda-rpc/internal/app/healthcheck"
	"github.com/alexliesenfeld/health"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"log/slog"
	"time"
)

func Register(router *gin.Engine, logger *slog.Logger) {
	global := router.Group("/api")
	global.Use(sloggin.New(logger.With("module", "api/global")))

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
	)

	global.GET("/health", gin.WrapF(health.NewHandler(checker)))
}
