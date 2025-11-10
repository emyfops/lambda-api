package middlewares

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

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
