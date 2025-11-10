package main

import (
	"crypto/rand"
	"crypto/rsa"
	"log/slog"
	"os"
	"runtime/debug"
	"time"

	"github.com/Edouard127/lambda-api/api"
	"github.com/Edouard127/lambda-api/api/middlewares"
	"github.com/Edouard127/lambda-api/internal"
	"github.com/gofiber/fiber/v2"
	fiblog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	flag "github.com/spf13/pflag"
)

var (
	isOnline      = flag.Bool("online", true, "Online-mode authentication")
	isDebug       = flag.Bool("debug", true, "Enable debug log output")
	redisEndpoint = flag.String("redis", "", "Endpoint of the standalone redis instance")
	keyPath       = flag.String("key", "", "Path to the private key")
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

	var key *rsa.PrivateKey
	var err error

	f, err := os.ReadFile(*keyPath)
	if err != nil {
		key, _ = rsa.GenerateKey(rand.Reader, 2048)
		logger.Warn("Failed to read the content of the private key", err)
	}

	key, err = jwt.ParseRSAPrivateKeyFromPEM(f)
	if err != nil {
		key, _ = rsa.GenerateKey(rand.Reader, 2048)
		logger.Warn("Failed to read the content of the private key", err)
	}

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
	internal.Set("key", key)

	if *isDebug {
		router.Use(fiblog.New(fiblog.Config{
			Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		}))
		router.Use(pprof.New())
	}

	router.Use(recover.New(recover.Config{EnableStackTrace: *isDebug}))

	api.New(router, rdb)

	panic(router.Listen(":8080"))
}

func printBuildInfo(logger *slog.Logger) {
	binaryInfo, _ := debug.ReadBuildInfo()
	settings := make(map[string]string)
	for _, v := range binaryInfo.Settings {
		settings[v.Key] = v.Value
	}
	logger.Debug("Build info", slog.Any("settings", settings))
}
