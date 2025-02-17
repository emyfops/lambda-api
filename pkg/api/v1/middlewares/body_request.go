package middlewares

import (
	"github.com/Edouard127/lambda-api/pkg/api/v1/models/response"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
)

var failedBodies = promauto.NewCounter(prometheus.CounterOpts{
	Name: "lambda_api_failed_logins",
	Help: "Total number of failed request with a defined body",
})

func BodyRequest[T any](ctx *gin.Context) {
	var body T

	err := ctx.Bind(&body)
	if err != nil {
		failedBodies.Inc()

		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ValidationError{
			Message: "Required fields are missing or invalid",
			Errors:  err.Error(),
		})
		return
	}

	ctx.Set("body", body)
	ctx.Next()
}
