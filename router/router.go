package router

import (
	"github.com/Edouard127/lambda-rpc/pkg/cache"
	v1 "github.com/Edouard127/lambda-rpc/router/api/v1"
	"github.com/gin-gonic/gin"
	"time"
)

func CreateEngine() *gin.Engine {
	router := gin.New()
	router.Use(cache.NewRateLimiter(time.Second*5, 2).Handle) // Rate limit requests (2 requests per 5 seconds)
	router.Use(gin.Recovery())                                // Recover from any panics

	// In the future we can add v2, v3, etc.
	apiv1 := router.Group("/api/v1")
	apiv1.Use()
	{
		apiv1.POST("/party", v1.CreateParty)
	}

	return router
}
