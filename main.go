//go:generate swag i -g main.go -dir ./pkg/api/v1/ --instanceName v1 -o openapi-spec

package main

import (
	"context"
	"flag"
	"github.com/Edouard127/lambda-api/cmd"
	_ "github.com/Edouard127/lambda-api/openapi-spec" // Required for swagger documentation
	"github.com/Edouard127/lambda-api/pkg/api/global"
	"github.com/Edouard127/lambda-api/pkg/api/global/middlewares"
	v1 "github.com/Edouard127/lambda-api/pkg/api/v1"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"net/http"
	"runtime/debug"
)

var isDebug = flag.Bool("debug", true, "Enable debug log output")

func main() {
	flag.Parse()
	var logger *zap.Logger
	if *isDebug {
		logger = zap.Must(zap.NewDevelopment())
	} else {
		logger = zap.Must(zap.NewProduction())
	}

	printBuildInfo(logger)

	// Create and initialize the redis connection
	rdb := redis.NewClient(cmd.RedisOptions())
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	gin.SetMode("") // read the mode from the test.v flag
	router := gin.New()
	router.SetTrustedProxies(nil)

	// Setup metrics
	router.Use(middlewares.PrometheusMiddleware())
	go startPrometheus(logger)

	// Prevent panics from crashing the server
	router.Use(gin.Recovery())

	// Register the APIs
	global.Register(rdb, router)

	v1.Register(rdb, router)
	router.GET("/swagger/v1/*any", ginSwagger.WrapHandler(swaggerFiles.NewHandler(), ginSwagger.InstanceName("v1")))

	// Return OK for the root path (helm chart test)
	router.GET("/", func(ctx *gin.Context) { ctx.String(http.StatusOK, "OK") })

	err = router.Run(":80")
	if err != nil {
		logger.Fatal("Server listening error", zap.Error(err))
	}
}

func startPrometheus(logger *zap.Logger) {
	logger.Info("Starting prometheus metrics server on :2112")

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9100", nil)
}

// printBuildInfo reading compile information of the binary program with runtime/debug package，and print it to log
func printBuildInfo(logger *zap.Logger) {
	binaryInfo, _ := debug.ReadBuildInfo()
	settings := make(map[string]string)
	for _, v := range binaryInfo.Settings {
		settings[v.Key] = v.Value
	}
	logger.Debug("Build info", zap.Any("settings", settings))
}
