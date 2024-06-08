//go:generate swag i -g main.go -dir .\pkg\api\v1\ --instanceName v1 -o openapi-spec

package main

import (
	"fmt"
	"github.com/Edouard127/lambda-rpc/internal/app/state"
	_ "github.com/Edouard127/lambda-rpc/openapi-spec"
	"github.com/Edouard127/lambda-rpc/pkg/api/global/middlewares"
	v1 "github.com/Edouard127/lambda-rpc/pkg/api/v1"
	"github.com/alexflint/go-arg"
	"github.com/gin-gonic/gin"
	"github.com/khaaleoo/gin-rate-limiter/core"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
	"log/slog"
	"os"
	"time"
)

// @Title Lambda RPC API
// @Version 1.0
// @Description This is the API for the Lambda Discord RPC handler
// @Contact.Name Lambda Discord
// @Contact.Url https://discord.gg/J23U4YEaAr
//
// @license.name GNU General Public License v3.0
// @license.url https://www.gnu.org/licenses/gpl-3.0.html
func main() {
	// Set the environment
	gin.SetMode(state.CurrentArgs.Environment)

	// Create a new logger
	logger := slog.New(slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: state.CurrentArgs.Verbose,
		}),
	)

	// Create a new router
	router := gin.New()

	// Setup metrics
	router.Use(middlewares.PrometheusMiddleware())
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Apply rate limiter after prometheus
	router.Use(gin.Recovery(), core.RequireRateLimiter(core.RateLimiter{
		RateLimiterType: core.IPRateLimiter,
		Key:             "iplimiter_maximum_requests_for_ip",
		Option: core.RateLimiterOption{
			Limit: rate.Limit(state.CurrentArgs.RateLimit),
			Burst: state.CurrentArgs.RateBurst,
			Len:   time.Duration(state.CurrentArgs.RatePunish) * time.Second,
		},
	}))

	// Provide swagger documentation
	router.GET("/swagger/v1/*any", ginSwagger.WrapHandler(swaggerFiles.NewHandler(), ginSwagger.InstanceName("v1")))

	// Register the APIs
	v1.Register(router, logger)
	// v2.Register(router, logger)
	// ...

	_ = router.Run(fmt.Sprintf(":%d", state.CurrentArgs.Port))
}

func init() {
	arg.MustParse(&state.CurrentArgs)
}
