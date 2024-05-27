package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

var (
	prometheusRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "lambda_rpc_requests_total",
		Help: "Total number of requests",
	}, []string{"path", "method", "status"})

	prometheusRequestsDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "lambda_rpc_requests_duration",
		Help:    "Duration of requests",
		Buckets: prometheus.DefBuckets,
	}, []string{"path", "method", "status"})
)

func PrometheusMiddleware() gin.HandlerFunc {
	prometheus.MustRegister(prometheusRequestsTotal)
	prometheus.MustRegister(prometheusRequestsDuration)

	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		duration := time.Since(start)
		status := ctx.Writer.Status()
		path := ctx.Request.URL.Path
		method := ctx.Request.Method

		prometheusRequestsTotal.WithLabelValues(path, method, strconv.Itoa(status)).Inc()
		prometheusRequestsDuration.WithLabelValues(path, method, strconv.Itoa(status)).Observe(duration.Seconds())
	}
}
