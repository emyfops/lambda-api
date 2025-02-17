package v1

import (
	"github.com/Edouard127/lambda-api/internal"
	"github.com/Edouard127/lambda-api/pkg/api/v1/endpoints"
	"github.com/Edouard127/lambda-api/pkg/api/v1/middlewares"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/request"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
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
func Register(cache *memcache.Client, router *gin.Engine) {
	v1 := router.Group("/api/v1")

	// Login endpoints
	v1.POST("/login", middlewares.BodyRequest[request.Authentication], endpoints.Login)
	v1.POST("/link/discord", middlewares.CheckAuth, endpoints.LinkDiscord)

	// Party endpoints
	v1.POST("/party/create", middlewares.CheckAuth, internal.With(cache, endpoints.CreateParty))
	v1.PUT("/party/join", middlewares.CheckAuth, middlewares.BodyRequest[request.JoinParty], internal.With(cache, endpoints.JoinParty))
	v1.PUT("/party/leave", middlewares.CheckAuth, internal.With(cache, endpoints.LeaveParty))
	v1.DELETE("/party/delete", middlewares.CheckAuth, internal.With(cache, endpoints.DeleteParty))
	v1.GET("/party", middlewares.CheckAuth, internal.With(cache, endpoints.GetParty))
}
