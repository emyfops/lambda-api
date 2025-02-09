package internal

import "github.com/gin-gonic/gin"

// MustGet returns the value of the key in the gin context.
func MustGet[T any](c *gin.Context, key string) T {
	return c.MustGet(key).(T)
}
