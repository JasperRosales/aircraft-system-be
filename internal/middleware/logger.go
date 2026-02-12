package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/JasperRosales/aircraft-system-be/internal/util"
)

func LoggerMiddleware(logger *util.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		logger.Info("Incoming request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"client_ip", c.ClientIP(),
			"latency", duration,
		)
	}
}
