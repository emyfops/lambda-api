package endpoints

import (
	"github.com/Edouard127/lambda-api/internal"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/request"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	"time"
)

var (
	successfulLogins = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "lambda_api_successful_logins",
		Help: "Total number of successful logins",
	}, []string{"version"})

	failedLogins = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "lambda_api_failed_logins",
		Help: "Total number of failed logins",
	}, []string{"version"})
)

// Login godoc
//
//	@Summary		Login to the server
//	@Description	Login to the server using a Minecraft username and a Mojang session hash
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			login	body		request.Authentication	true	"Authentication"
//	@Success		200		{object}	response.Authentication
//	@Failure		400		{object}	response.ValidationError
//	@Failure		401		{object}	response.Error
//	@Failure		500		{object}	response.Error
//	@Router			/login 	[post]
func Login(ctx *gin.Context) {
	login := ctx.MustGet("body").(request.Authentication)

	player, err := response.GetPlayer(login.Username, login.Hash)
	if err != nil {
		failedLogins.WithLabelValues("v1").Inc()
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Error{
			Message: "Invalid credentials",
		})
		return
	}

	signed, err := internal.NewJwt(player)
	if err != nil {
		failedLogins.WithLabelValues("v1").Inc()
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Error{
			Message: "Failed to create token",
		})
		return
	}

	successfulLogins.WithLabelValues("v1").Inc()

	ctx.AbortWithStatusJSON(http.StatusOK, response.Authentication{
		AccessToken: signed,
		ExpiresIn:   int64(time.Hour * 24),
		TokenType:   "Bearer",
	})
}
