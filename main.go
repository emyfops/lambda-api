//go:generate swag i -g main.go -dir .\pkg\api\v1\ --instanceName v1 -o api/openapi-spec

package main

import (
	"bytes"
	"fmt"
	_ "github.com/Edouard127/lambda-rpc/api/openapi-spec"
	v1 "github.com/Edouard127/lambda-rpc/pkg/api/v1"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io"
)

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	router := gin.New()

	router.GET("/swagger/v1/*any", ginSwagger.WrapHandler(swaggerFiles.NewHandler(), ginSwagger.InstanceName("v1")))
	router.Use(gin.Recovery())
	router.Use(DebugMiddleware())

	v1.Register(router) // Register the v1 API

	_ = router.Run(":8080")
}

func DebugMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var buf bytes.Buffer
		tee := io.TeeReader(ctx.Request.Body, &buf)
		body, _ := io.ReadAll(tee)
		ctx.Request.Body = io.NopCloser(&buf)
		fmt.Println(
			fmt.Sprintf(
				"Method: %s\nURL: %s\nHeaders: %s",
				ctx.Request.Method,
				ctx.Request.URL,
				body),
		)
		ctx.Next()
	}
}
