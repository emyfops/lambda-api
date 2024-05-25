package auth

import "github.com/gin-gonic/gin"

// GinMustGet returns the value of the key in the gin context.
// If the key does not exist, it panics.
func GinMustGet[T any](c *gin.Context, key string) T {
	v := c.MustGet(key)
	return v.(T)
}
