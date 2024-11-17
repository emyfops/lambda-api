//go:generate swag i -g main.go -dir ./pkg/api/v1/ --instanceName v1 -o openapi-spec

package main

import (
	"github.com/Edouard127/lambda-rpc/internal/app/state"
	_ "github.com/Edouard127/lambda-rpc/openapi-spec"
	"github.com/Edouard127/lambda-rpc/pkg/api/global"
	"github.com/Edouard127/lambda-rpc/pkg/api/global/middlewares"
	v1 "github.com/Edouard127/lambda-rpc/pkg/api/v1"
	"github.com/alexflint/go-arg"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log/slog"
	"net/http"
	"os"
)

// @Title Lambda API
// @Version 1.0
// @Description This is the official API for Lambda Client
// @Contact.Name Lambda Discord
// @Contact.Url https://discord.gg/QjfBxJzE5x
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
			Level: state.CurrentArgs.LogLevel,
		}),
	)

	// Create a new router
	router := gin.New()
	router.SetTrustedProxies(nil)

	// Setup metrics
	router.Use(middlewares.PrometheusMiddleware())
	go http.ListenAndServe(":9100", promhttp.Handler())

	// Provide swagger documentation
	router.GET("/swagger/v1/*any", ginSwagger.WrapHandler(swaggerFiles.NewHandler(), ginSwagger.InstanceName("v1")))

	// Prevent panics from crashing the server
	router.Use(gin.Recovery())

	// Register the APIs
	global.Register(router, logger)
	v1.Register(router, logger)
	// v2.Register(router, logger)
	// ...

	// Return OK for the root path (helm chart test)
	router.GET("/", func(ctx *gin.Context) { ctx.String(http.StatusNoContent, "OK") })

	_ = router.Run(":80")
}

func init() {
	arg.MustParse(&state.CurrentArgs)
}
