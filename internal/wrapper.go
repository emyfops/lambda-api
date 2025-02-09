package internal

import (
	"github.com/gin-gonic/gin"
)

// With wraps a gin.HandlerFunc call with a value of T
func With[T any](v T, fn func(*gin.Context, T)) gin.HandlerFunc {
	return func(ctx *gin.Context) { fn(ctx, v) }
}
