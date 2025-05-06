package middlewares

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	"strconv"
	"time"
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Subsystem: "http",
	Name:      "request_duration_seconds",
	Help:      "Duration of HTTP requests in seconds",
	Buckets:   prometheus.DefBuckets,
}, []string{"path", "method", "status"})

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	code := http.StatusInternalServerError
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	} else {
		e = &fiber.Error{
			Code:    code,
			Message: err.Error(),
		}
	}
	return ctx.Status(code).JSON(e)
}

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

func RequestDuration() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		s := time.Now()

		err := ctx.Next()

		duration := time.Since(s).Seconds()
		httpDuration.WithLabelValues(ctx.Path(), ctx.Method(), strconv.Itoa(ctx.Response().StatusCode())).Observe(duration)

		return err
	}
}
