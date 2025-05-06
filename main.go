package main

import (
	"github.com/Edouard127/lambda-api/api"
	"github.com/Edouard127/lambda-api/api/middlewares"
	"github.com/Edouard127/lambda-api/internal"
	"github.com/gofiber/fiber/v2"
	fiblog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	flag "github.com/spf13/pflag"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"time"
)

var (
	isOnline      = flag.Bool("online", true, "Online-mode authentication")
	isDebug       = flag.Bool("debug", true, "Enable debug log output")
	redisEndpoint = flag.String("redis", "", "Endpoint of the standalone redis instance")
)

func main() {
	flag.Parse()

	var logger *slog.Logger
	if *isDebug {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	printBuildInfo(logger)
	go startPrometheus(logger)

	if !*isOnline {
		logger.Warn("Warning, running in offline mode allows users to spoof their authentication and usurpate other players")
	}

	rdb := redis.NewClient(&redis.Options{Addr: *redisEndpoint, ReadTimeout: 1 * time.Second})

	router := fiber.New(fiber.Config{
		Network:      "tcp", // v4 and v6
		ReadTimeout:  5 * time.Second,
		ErrorHandler: middlewares.ErrorHandler,
	})

	internal.Set("logger", logger)
	internal.Set("cache", rdb)

	if *isDebug {
		router.Use(fiblog.New(fiblog.Config{
			Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		}))
		router.Use(pprof.New())
	}
	router.Use(middlewares.RequestDuration())
	router.Use(recover.New(recover.Config{EnableStackTrace: *isDebug}))

	api.New(router, rdb)

	panic(router.Listen(":8080"))
}

func startPrometheus(logger *slog.Logger) {
	logger.Info("Starting prometheus metrics server on :9100")

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9100", nil)
}

func printBuildInfo(logger *slog.Logger) {
	binaryInfo, _ := debug.ReadBuildInfo()
	settings := make(map[string]string)
	for _, v := range binaryInfo.Settings {
		settings[v.Key] = v.Value
	}
	logger.Debug("Build info", slog.Any("settings", settings))
}
