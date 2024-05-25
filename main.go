//go:generate swag i -g main.go -dir .\pkg\api\v1\ --instanceName v1 -o openapi-spec

package main

import (
	"fmt"
	"github.com/Edouard127/lambda-rpc/internal/app/state"
	_ "github.com/Edouard127/lambda-rpc/openapi-spec"
	v1 "github.com/Edouard127/lambda-rpc/pkg/api/v1"
	"github.com/alexflint/go-arg"
	"github.com/gin-gonic/gin"
	"github.com/khaaleoo/gin-rate-limiter/core"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
	"log/slog"
	"os"
	"time"
)

var _ = arg.MustParse(&state.CurrentArgs)
var limiter = core.RateLimiter{
	RateLimiterType: core.IPRateLimiter,
	Key:             "iplimiter_maximum_requests_for_ip",
	Option: core.RateLimiterOption{
		Limit: rate.Limit(state.CurrentArgs.RateLimit),
		Burst: state.CurrentArgs.RateBurst,
		Len:   time.Duration(state.CurrentArgs.RatePunish) * time.Second,
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
	gin.SetMode(state.CurrentArgs.Environment)

	router := gin.New()
	__prometheus(router)

	logger := slog.New(slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: state.CurrentArgs.Verbose,
		}),
	)

	router.Use(gin.Recovery())
	router.Use(core.RequireRateLimiter(limiter))

	router.GET("/swagger/v1/*any", ginSwagger.WrapHandler(swaggerFiles.NewHandler(), ginSwagger.InstanceName("v1")))

	v1.Register(router, logger)

	_ = router.Run(fmt.Sprintf(":%d", state.CurrentArgs.Port))
}

func __prometheus(router *gin.Engine) {
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewGoCollector())
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewBuildInfoCollector())

	// TODO: Should I protected this endpoint ?
	router.GET("/metrics", gin.WrapH(promhttp.HandlerFor(reg, promhttp.HandlerOpts{})))
}
