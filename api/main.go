package api

import (
	"github.com/Edouard127/lambda-api/api/middlewares"
	"github.com/Edouard127/lambda-api/api/routes"
	"github.com/Edouard127/lambda-api/internal"
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
func Register(router *gin.Engine, cache memcached.Client) {
	v1 := router.Group("/api/v1")

	// Login endpoints
	v1.POST("/login", routes.Login)
	v1.POST("/link/discord", middlewares.CheckAuth, internal.With(cache, routes.LinkDiscord))

	// Party endpoints
	v1.POST("/party/create", middlewares.CheckAuth, middlewares.DiscordCheck, internal.With(cache, routes.CreateParty))
	v1.PUT("/party/join", middlewares.CheckAuth, middlewares.DiscordCheck, internal.With(cache, routes.JoinParty))
	v1.PUT("/party/leave", middlewares.CheckAuth, middlewares.DiscordCheck, internal.With(cache, routes.LeaveParty))
	v1.DELETE("/party/delete", middlewares.CheckAuth, middlewares.DiscordCheck, internal.With(cache, routes.DeleteParty))
	v1.GET("/party", middlewares.CheckAuth, middlewares.DiscordCheck, internal.With(cache, routes.GetParty))
	v1.GET("/party/listen", middlewares.CheckAuth, middlewares.DiscordCheck, internal.With(cache, routes.PartyListen))

	// Cape endpoints
	v1.GET("/cape", internal.With(cache, routes.GetCape))
	v1.PUT("/cape", middlewares.CheckAuth, internal.With(cache, routes.SetCape))
}
