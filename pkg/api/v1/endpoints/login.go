package endpoints

import (
	"github.com/Edouard127/lambda-api/internal"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/request"
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
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

// Login allows a player to log in to the server using a Minecraft username and Mojang session hash.
//
//	@Summary	Login to the server
//	@Tags		Authentication
//	@Accept		json
//	@Produce	json
//	@Param		login	body	request.Authentication	true	"Authentication credentials (Minecraft username and Mojang session hash)"
//	@Success	200	{object}	response.Authentication		"Successfully logged in and retrieved authentication token"
//	@Failure	400	{object}	response.ValidationError	"Invalid or missing authentication fields"
//	@Failure	401	{object}	response.Error				"Invalid credentials"
//	@Failure	500	{object}	response.Error				"Internal server error"
//	@Router		/login [post]
func Login(ctx *gin.Context) {
	logger := ctx.MustGet("logger").(*zap.Logger)

	var login request.Authentication

	err := ctx.Bind(&login)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ValidationError{
			Message: "Required fields are missing or invalid",
			Errors:  err.Error(),
		})
		return
	}

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
		logger.Error("Error signing token", zap.Error(err))
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
