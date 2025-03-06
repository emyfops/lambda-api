package v1

import (
	"github.com/Edouard127/lambda-api/internal"
	"github.com/Edouard127/lambda-api/pkg/api/v1/endpoints"
	"github.com/Edouard127/lambda-api/pkg/api/v1/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/yeqown/memcached"
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
func Register(router *gin.Engine, cache memcached.Client, persistent memcached.Client) {
	v1 := router.Group("/api/v1")

	// Login endpoints
	v1.POST("/login", endpoints.Login)
	v1.POST("/link/discord", middlewares.CheckAuth, internal.With(cache, endpoints.LinkDiscord))

	// Party endpoints
	v1.POST("/party/create", middlewares.CheckAuth, middlewares.DiscordCheck, internal.With(cache, endpoints.CreateParty))
	v1.PUT("/party/join", middlewares.CheckAuth, middlewares.DiscordCheck, internal.With(cache, endpoints.JoinParty))
	v1.PUT("/party/leave", middlewares.CheckAuth, middlewares.DiscordCheck, internal.With(cache, endpoints.LeaveParty))
	v1.DELETE("/party/delete", middlewares.CheckAuth, middlewares.DiscordCheck, internal.With(cache, endpoints.DeleteParty))
	v1.GET("/party", middlewares.CheckAuth, middlewares.DiscordCheck, internal.With(cache, endpoints.GetParty))
	v1.GET("/party/listen", middlewares.CheckAuth, middlewares.DiscordCheck, internal.With(cache, endpoints.PartyListen))

	// Cape endpoints
	v1.GET("/cape", internal.With(persistent, endpoints.GetCape))
	v1.PUT("/cape", middlewares.CheckAuth, internal.With(persistent, endpoints.SetCape))
}
