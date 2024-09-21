package global

import (
	"github.com/Edouard127/lambda-rpc/internal/app/healthz"
	"github.com/gin-gonic/gin"
	"github.com/heptiolabs/healthcheck"
	sloggin "github.com/samber/slog-gin"
	"log/slog"
	"time"
)

func Register(router *gin.Engine, logger *slog.Logger) {
	global := router.Group("/api")
	global.Use(sloggin.New(logger.With("module", "api/global")))

	// Initialize the healthcheck handler
	health := healthcheck.NewHandler()

	// Liveness checks
	// health.AddLivenessCheck(...)

	// Readiness checks
	health.AddReadinessCheck("http-connection-mojang-session",
		healthcheck.Async(
			healthz.HTTPGetCheck("https://sessionserver.mojang.com/session/minecraft/hasJoined", 2000*time.Millisecond), 60*time.Second),
	)

	global.GET("/live", gin.WrapF(health.LiveEndpoint))
	global.GET("/ready", gin.WrapF(health.ReadyEndpoint))
}
