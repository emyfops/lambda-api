package middlewares

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	"strconv"
)

var (
	cfg = metrics.Config{
		DurationBuckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
		MethodLabel:     "method",
		StatusCodeLabel: "code",
	}

	handler = middleware.New(middleware.Config{
		Recorder:      metrics.NewRecorder(cfg),
		GroupedStatus: true,
		IgnoredPaths:  []string{"/api/health", "/favicon.ico", "/robots.txt"},
	})

	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: cfg.Prefix,
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "The number of requests on a given router.",
	}, []string{cfg.MethodLabel, cfg.StatusCodeLabel})
)

func Locals(args ...any) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		for i := 0; i < len(args); i += 2 {
			if i+1 < len(args) {
				ctx.Locals(args[i], args[i+1])
			}
		}
		return ctx.Next()
	}
}

func MeasureRequest() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		handler.Measure("", reporter{ctx}, func() {})

		err := ctx.Next()

		httpRequestsTotal.WithLabelValues(string(ctx.Request().Header.Method()), strconv.Itoa(ctx.Response().StatusCode()))

		return err
	}
}

type reporter struct {
	c *fiber.Ctx
}

func (r reporter) Method() string {
	return r.c.Method()
}

func (r reporter) Context() context.Context {
	return r.c.Context()
}

func (r reporter) URLPath() string {
	return r.c.Path()
}

func (r reporter) StatusCode() int {
	return r.c.Response().StatusCode()
}

func (r reporter) BytesWritten() int64 {
	return int64(len(r.c.Response().Body()))
}
