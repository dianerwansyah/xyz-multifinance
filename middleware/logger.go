package middleware

import (
	"time"

	"xyz-multifinance/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		statusCode := c.Writer.Status()
		duration := time.Since(start)

		entry := logger.Log.WithFields(logrus.Fields{
			"method":     method,
			"path":       path,
			"status":     statusCode,
			"duration":   duration,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		})

		if len(c.Errors) > 0 {
			entry = entry.WithField("error", c.Errors.String())
		}

		if statusCode >= 500 {
			entry.Error("incoming request")
		} else if statusCode >= 400 {
			entry.Warn("incoming request")
		} else {
			entry.Info("incoming request")
		}
	}
}
