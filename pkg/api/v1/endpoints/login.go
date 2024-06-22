package endpoints

import (
	"github.com/Edouard127/lambda-rpc/internal/app/auth"
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models/request"
	"github.com/Edouard127/lambda-rpc/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

var (
	successfulLogins = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "lambda_rpc_successful_logins",
		Help: "Total number of successful logins",
	}, []string{"version"})

	failedLogins = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "lambda_rpc_failed_logins",
		Help: "Total number of failed logins",
	}, []string{"version"})
)

// Login godoc
// @BasePath /api/v1
// @Summary Login to the server
// @Description Login to the server using a Discord identify token, a Minecraft username and a Mojang session hash
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login body request.Authentication true "Authentication"
// @Success 200 {object} response.Authentication
// @Failure 400 {object} response.ValidationError
// @Failure 401 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /party/login [post]
func Login(ctx *gin.Context) {
	var login request.Authentication
	if err := ctx.Bind(&login); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ValidationError{
			Message: "Required fields are missing or invalid",
			Errors:  err.Error(),
		})
		failedLogins.WithLabelValues("v1").Inc()
	}

	player, err := response.GetPlayer(login.Token, login.Username, login.Hash)
	if err != nil {
		failedLogins.WithLabelValues("v1").Inc()
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Error{
			Message: "Invalid credentials",
		})
		failedLogins.WithLabelValues("v1").Inc()
	}

	signed, err := auth.CreateJwtToken(player)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{
			Message: "Failed to create token",
		})
		failedLogins.WithLabelValues("v1").Inc()
	}

	ctx.AbortWithStatusJSON(http.StatusOK, response.Authentication{
		AccessToken: signed,
		ExpiresIn:   int64(time.Hour * 24),
		TokenType:   "Bearer",
	})

	successfulLogins.WithLabelValues("v1").Inc()
}
