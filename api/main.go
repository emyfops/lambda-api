package api

import (
	"context"
	"github.com/Edouard127/lambda-api/api/middlewares"
	"github.com/Edouard127/lambda-api/api/routes"
	"github.com/Edouard127/lambda-api/internal"
	"github.com/alexliesenfeld/health"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/redis/go-redis/v9"
	"time"
)

func New(router fiber.Router, cache *redis.Client) {
	api := router.Group("/api")

	checker := health.NewChecker(
		health.WithCacheDuration(2*time.Second),
		health.WithTimeout(5*time.Second),
		health.WithPeriodicCheck(
			2*time.Second,
			time.Second,
			health.Check{
				Name:  "http-connection-mojang-session",
				Check: internal.HTTPGetCheck("https://sessionserver.mojang.com/session/minecraft/hasJoined"),
			},
		),
		health.WithPeriodicCheck(
			60*time.Second,
			time.Second,
			health.Check{
				Name:  "redis-connection",
				Check: func(ctx context.Context) error { return cache.Ping(ctx).Err() },
			},
		),
	)
	api.Get("/health", adaptor.HTTPHandler(health.NewHandler(checker)))

	v1 := router.Group("/api/v1")

	// Login
	v1.Post("/login", routes.Login)
	v1.Post("/link/discord", middlewares.MinecraftCheck, routes.LinkDiscord)

	// Capes
	v1.Get("/cape", routes.GetCape)
	v1.Get("/capes", routes.GetCapes)
	v1.Put("/cape", middlewares.MinecraftCheck, routes.SetCape)

	// Party endpoints
	/*
		v1.Post("/party/create", middlewares.CheckAuth, middlewares.DiscordCheck, internal.Locals(cache, routes.CreateParty))
		v1.Put("/party/join", middlewares.CheckAuth, middlewares.DiscordCheck, internal.Locals(cache, routes.JoinParty))
		v1.Put("/party/leave", middlewares.CheckAuth, middlewares.DiscordCheck, internal.Locals(cache, routes.LeaveParty))
		v1.Delete("/party/delete", middlewares.CheckAuth, middlewares.DiscordCheck, internal.Locals(cache, routes.DeleteParty))
		v1.Get("/party", middlewares.CheckAuth, middlewares.DiscordCheck, internal.Locals(cache, routes.GetParty))
		v1.Get("/party/listen", middlewares.CheckAuth, middlewares.DiscordCheck, internal.Locals(cache, routes.PartyListen))
	*/
}
