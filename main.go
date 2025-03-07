//go:generate swag i -g main.go -dir ./pkg/api/v1/ --instanceName v1 -o api

package main

import (
	"github.com/Edouard127/lambda-api/pkg/api/global"
	"github.com/Edouard127/lambda-api/pkg/api/global/middlewares"
	v1 "github.com/Edouard127/lambda-api/pkg/api/v1"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	flag "github.com/spf13/pflag"
	"github.com/yeqown/memcached"
	"go.uber.org/zap"
	"net/http"
	"runtime/debug"
	"strings"
)

func main() {
	var isDebug bool
	var staging string
	var dragons []string

	flag.Bool("online", true, "Online-mode")
	flag.StringVar(&staging, "staging", "debug", "Gin staging mode (info, debug, release)")
	flag.BoolVar(&isDebug, "debug", true, "Enable debug log output")
	flag.StringArrayVar(&dragons, "nodes", []string{}, "Memcache nodes")

	flag.Parse()

	var logger *zap.Logger
	if isDebug {
		logger = zap.Must(zap.NewDevelopment())
	} else {
		logger = zap.Must(zap.NewProduction())
	}

	printBuildInfo(logger)

	dragon, err := memcached.New(strings.Join(dragons, ","))
	if err != nil {
		logger.Fatal("Failed to connect to DragonFly", zap.Error(err))
	}

	gin.SetMode(staging)
	router := gin.New()
	router.SetTrustedProxies(nil)

	go startPrometheus(logger)

	if isDebug {
		router.Use(gin.Logger())
	}
	router.Use(middlewares.PrometheusMiddleware())
	router.Use(gin.Recovery())
	router.Use(middlewares.Logger(logger))

	global.Register(router)
	v1.Register(router, dragon)

	err = router.Run(":8080")
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
