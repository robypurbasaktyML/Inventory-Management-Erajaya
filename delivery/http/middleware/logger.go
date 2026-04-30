package middleware

import (
	"time"

	"Inventory-Management-Erajaya/pkg/logger"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		entry := logger.Log.WithFields(map[string]interface{}{
			"status":  status,
			"method":  c.Request.Method,
			"path":    path,
			"latency": latency.String(),
			"ip":      c.ClientIP(),
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.String())
		} else if status >= 500 {
			entry.Error("server error")
		} else if status >= 400 {
			entry.Warn("client error")
		} else {
			entry.Info("request")
		}
	}
}
