package main

import (
	_ "github.com/Edouard127/lambda-rpc/api/openapi-spec"
	"github.com/Edouard127/lambda-rpc/pkg/api/v1"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	router := gin.New()

	router.GET("/swagger/v1/*any", ginSwagger.WrapHandler(swaggerFiles.NewHandler(), ginSwagger.InstanceName("v1")))
	router.Use(gin.Recovery())

	v1.Register(router) // Register the v1 API

	_ = router.Run(":8080")
}
