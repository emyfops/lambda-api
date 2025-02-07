package v1

import (
	"github.com/Edouard127/lambda-api/internal/app/gonic"
	"github.com/Edouard127/lambda-api/pkg/api/v1/endpoints"
	"github.com/Edouard127/lambda-api/pkg/api/v1/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Register godoc
//
//	@Title						Lambda API
//	@Version					v1
//	@Description				This is the official API for Lambda Client
//
//	@BasePath					/api/v1
//	@Schemes					http https
//
//	@Contact.name				Lambda Discord
//	@Contact.url				https://discord.gg/QjfBxJzE5x
//
//	@License.name				GNU General Public License v3.0
//	@License.url				https://www.gnu.org/licenses/gpl-3.0.html
//
//	@SecurityDefinitions.ApiKey	Bearer
//	@In							header
//	@Name						Authorization
//	@Description				Type "Bearer" followed by a space and JWT token.
func Register(client *redis.Client, router *gin.Engine) {
	v1 := router.Group("/api/v1")

	// Login endpoints
	v1.POST("/login", endpoints.Login)

	// Party endpoints
	v1.POST("/party/create", middlewares.CheckAuth, gonic.With(client, endpoints.CreateParty))
	v1.PUT("/party/join", middlewares.CheckAuth, gonic.With(client, endpoints.JoinParty))
	v1.PATCH("/party/edit", middlewares.CheckAuth, gonic.With(client, endpoints.EditParty))
	v1.PUT("/party/leave", middlewares.CheckAuth, gonic.With(client, endpoints.LeaveParty))
	v1.DELETE("/party/delete", middlewares.CheckAuth, gonic.With(client, endpoints.DeleteParty))
	v1.GET("/party", middlewares.CheckAuth, gonic.With(client, endpoints.GetParty))
}
