package v1

import (
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/endpoints"
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/middlewares"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"log/slog"
)

//
// @BasePath /api/v1
//
// @SecurityDefinitions.Apikey Bearer
// @In header
// @Name Authorization
// @Description Type "Bearer" followed by a space and JWT token.

func Register(router *gin.Engine, logger *slog.Logger) {
	v1 := router.Group("/api/v1")
	v1.Use(sloggin.New(logger.With("module", "api/v1")))

	v1.POST("/login", endpoints.Login)
	v1.POST("/party/create", middlewares.CheckAuth, endpoints.CreateParty)
	v1.PUT("/party/join", middlewares.CheckAuth, endpoints.JoinParty)
	v1.PATCH("/party/edit", middlewares.CheckAuth, endpoints.EditParty)
}
