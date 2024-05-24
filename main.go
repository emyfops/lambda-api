//go:generate swag i -g main.go -dir .\pkg\api\v1\ --instanceName v1 -o openapi-spec

package main

import (
	"bytes"
	"fmt"
	_ "github.com/Edouard127/lambda-rpc/openapi-spec"
	v1 "github.com/Edouard127/lambda-rpc/pkg/api/v1"
	"github.com/alexflint/go-arg"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io"
	"net/netip"
)

var args struct {
	Host netip.Addr `arg:"-h,--host" help:"Host address, supports v4 and v6" default:"127.0.0.1"`
	Port int        `arg:"-p,--port" help:"Port number" default:"8080"`

	Verbose string `arg:"-v,--verbose" help:"Log level" default:"info" placeholder:"INFO | DEBUG | WARN | ERROR"`

	AllowInsecure bool `arg:"--allow-insecure" help:"Allow insecure minecraft accounts to connect" default:"false"`
}

var _ = arg.MustParse(&args)

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	fmt.Println(args)
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
