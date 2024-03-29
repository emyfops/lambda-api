package v1

import (
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/endpoints"
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/middlewares"
	"github.com/gin-gonic/gin"
)

// @Title Lambda RPC API
// @Version 1.0
// @Description This is the API for the Lambda Discord RPC handler
// @Contact.Name Lambda Discord
// @Contact.Url https://discord.gg/J23U4YEaAr
//
// @license.name GNU General Public License v3.0
// @license.url https://www.gnu.org/licenses/gpl-3.0.html
//
// @BasePath /api/v1
//
// @SecurityDefinitions.Apikey Bearer
// @In header
// @Name Authorization
// @Description Type "Bearer" followed by a space and JWT token.

func Register(router *gin.Engine) {
	v1 := router.Group("/api/v1")

	v1.POST("/party/login", endpoints.Login)
	v1.POST("/party/create", middlewares.CheckAuth, endpoints.CreateParty)
	v1.POST("/party/join", middlewares.CheckAuth, endpoints.JoinParty)
}
