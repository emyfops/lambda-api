package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	ginmiddleware "github.com/slok/go-http-metrics/middleware/gin"
)

var handler = middleware.New(middleware.Config{
	Recorder: prometheus.NewRecorder(prometheus.Config{
		DurationBuckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
	}),
	GroupedStatus: true,
	IgnoredPaths:  []string{"/api/health", "/favicon.ico", "/robots.txt"},
})

func Metrics() gin.HandlerFunc {
	return ginmiddleware.Handler("", handler)
}
