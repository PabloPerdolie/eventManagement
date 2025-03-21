package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		status := c.Writer.Status()

		logger.Infow("Request processed",
			"status", status,
			"method", method,
			"path", path,
			"latency", latency,
			"client_ip", c.ClientIP(),
		)
	}
}
