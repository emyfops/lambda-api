//go:generate swag i -g main.go -dir ./pkg/api/v1/ --instanceName v1 -o openapi-spec

package main

import (
	"flag"
	_ "github.com/Edouard127/lambda-api/openapi-spec" // Required for swagger documentation
	"github.com/Edouard127/lambda-api/pkg/api/global"
	"github.com/Edouard127/lambda-api/pkg/api/global/middlewares"
	v1 "github.com/Edouard127/lambda-api/pkg/api/v1"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"net/http"
	"runtime/debug"
)

var _ = flag.Bool("insecure", false, "Insecure login")
var mode = flag.String("staging", "debug", "Gin staging mode (info, debug, release)")
var isDebug = flag.Bool("debug", true, "Enable debug log output")
var cacheNodes []string // list of memcached instances

func main() {
	flag.Parse()
	cacheNodes = flag.Args()

	var logger *zap.Logger
	if *isDebug {
		logger = zap.Must(zap.NewDevelopment())
	} else {
		logger = zap.Must(zap.NewProduction())
	}

	printBuildInfo(logger)

	mc := memcache.New(cacheNodes...)
	if err := mc.Ping(); err != nil {
		panic(err)
	}

	gin.SetMode(*mode)
	router := gin.New()
	router.SetTrustedProxies(nil)

	// Setup metrics
	router.Use(middlewares.PrometheusMiddleware())
	go startPrometheus(logger)

	// Prevent panics from crashing the server
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Register the APIs
	global.Register(router)

	v1.Register(mc, router)
	router.GET("/swagger/v1/*any", ginSwagger.WrapHandler(swaggerFiles.NewHandler(), ginSwagger.InstanceName("v1")))

	err := router.Run(":8080")
	if err != nil {
		logger.Fatal("Server listening error", zap.Error(err))
	}
}

func startPrometheus(logger *zap.Logger) {
	logger.Info("Starting prometheus metrics server on :9100")

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
