package middleware

import (
	"net/http"
	"runtime/debug"

	"xyz-multifinance/logger"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Log.WithFields(map[string]interface{}{
					"error": rec,
					"stack": string(debug.Stack()),
				}).Error("panic recovered")

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		}()

		c.Next()
	}
}
