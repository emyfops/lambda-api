package healthcheck

import (
	"context"
	"github.com/Edouard127/lambda-api/internal"
	"github.com/alexliesenfeld/health"
	"github.com/gin-gonic/gin"
	"github.com/yeqown/memcached"
	"time"
)

func Register(router *gin.Engine, cache memcached.Client) {
	global := router.Group("/api")

	// Initialize the healthcheck handler
	checker := health.NewChecker(
		health.WithCacheDuration(2*time.Second),
		health.WithTimeout(5*time.Second),
		health.WithPeriodicCheck(
			2*time.Second,
			time.Second,
			health.Check{
				Name:  "http-connection-mojang-session",
				Check: internal.HTTPGetCheck("https://sessionserver.mojang.com/session/minecraft/hasJoined"),
			},
		),
		health.WithPeriodicCheck(
			60*time.Second,
			time.Second,
			health.Check{
				Name:  "memcache-connection",
				Check: func(ctx context.Context) (err error) { _, err = cache.Version(ctx); return },
			},
		),
	)

	global.GET("/health", gin.WrapF(health.NewHandler(checker)))
}
