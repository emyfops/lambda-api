//go:generate swag i -g main.go -dir ./pkg/api/v1/ --instanceName v1 -o openapi-spec

package main

import (
	"context"
	"github.com/Edouard127/lambda-api/cmd"
	"github.com/Edouard127/lambda-api/pkg/api/global"
	"github.com/Edouard127/lambda-api/pkg/api/global/middlewares"
	v1 "github.com/Edouard127/lambda-api/pkg/api/v1"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/Edouard127/lambda-api/openapi-spec" // Required for swagger documentation
)

func main() {
	// Set the environment
	gin.SetMode(cmd.Arguments().Environment)

	// Create a new logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cmd.Arguments().LogLevel}))

	// Create and initialize the redis connection
	rdb := redis.NewClient(cmd.RedisOptions())
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	// Create a new router
	router := gin.New()
	router.SetTrustedProxies(nil)

	// Setup metrics
	router.Use(middlewares.PrometheusMiddleware())
	go func() {
		err := http.ListenAndServe(":9100", promhttp.Handler())
		if err != nil {
			logger.Error("Failed to start prometheus metrics", err)
		}
	}()

	// Prevent panics from crashing the server
	router.Use(gin.Recovery())

	// Register the APIs
	global.Register(rdb, router, logger)

	v1.Register(rdb, router, logger)
	router.GET("/swagger/v1/*any", ginSwagger.WrapHandler(swaggerFiles.NewHandler(), ginSwagger.InstanceName("v1")))

	// v2.Register(router, logger)
	// router.GET("/swagger/v2/*any", ginSwagger.WrapHandler(swaggerFiles.NewHandler(), ginSwagger.InstanceName("v2")))

	// Return OK for the root path (helm chart test)
	router.GET("/", func(ctx *gin.Context) { ctx.String(http.StatusOK, "OK") })

	err = router.Run(":8080")
}
