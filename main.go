//go:generate swag i -g main.go -dir .\pkg\api\v1\ --instanceName v1 -o openapi-spec

package main

import (
	"github.com/Edouard127/lambda-rpc/internal/app/state"
	_ "github.com/Edouard127/lambda-rpc/openapi-spec"
	v1 "github.com/Edouard127/lambda-rpc/pkg/api/v1"
	"github.com/alexflint/go-arg"
	"github.com/gin-gonic/gin"
	"github.com/khaaleoo/gin-rate-limiter/core"
	sloggin "github.com/samber/slog-gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log/slog"
	"os"
	"time"
)

var _ = arg.MustParse(&state.CurrentArgs)
var limiter = core.RateLimiter{
	RateLimiterType: core.IPRateLimiter,
	Key:             "iplimiter_maximum_requests_for_ip",
	Option: core.RateLimiterOption{
		Limit: 2,                // Maximum number of requests allowed per second
		Burst: 5,                // Maximum number of requests allowed to burst
		Len:   10 * time.Second, // Time window
	},
}

// @Title Lambda RPC API
// @Version 1.0
// @Description This is the API for the Lambda Discord RPC handler
// @Contact.Name Lambda Discord
// @Contact.Url https://discord.gg/J23U4YEaAr
//
// @license.name GNU General Public License v3.0
// @license.url https://www.gnu.org/licenses/gpl-3.0.html
func main() {
	router := gin.New()

	logger := slog.New(slog.NewJSONHandler(
		os.Stdout,
		nil),
	)
	router.Use(sloggin.New(logger), gin.Recovery())
	router.Use(core.RequireRateLimiter(limiter))

	router.GET("/swagger/v1/*any", ginSwagger.WrapHandler(swaggerFiles.NewHandler(), ginSwagger.InstanceName("v1")))

	v1.Register(router) // Register the v1 API

	_ = router.Run(":8080")
}
