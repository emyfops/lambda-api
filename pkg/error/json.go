package error

import (
	"github.com/gin-gonic/gin"
)

func JsonError(code Code) gin.H {
	return JsonResponse(code, code.String())
}

func JsonResponse(code Code, message string) gin.H {
	return gin.H{
		"code":    code,
		"message": message,
	}
}
