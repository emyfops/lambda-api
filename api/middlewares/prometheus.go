package middlewares

import (
	"github.com/Edouard127/lambda-api/api/metrics"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		duration := time.Since(start)
		status := ctx.Writer.Status()
		path := ctx.Request.URL.Path
		method := ctx.Request.Method

		metrics.RequestsTotal.WithLabelValues(path, method, strconv.Itoa(status)).Inc()
		metrics.RequestsDuration.WithLabelValues(path, method, strconv.Itoa(status)).Observe(duration.Seconds())
	}
}
