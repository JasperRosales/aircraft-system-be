package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/JasperRosales/aircraft-system-be/internal/util"
)

func LoggerMiddleware(logger *util.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		logger.Info("Incoming request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.String("client_ip", c.ClientIP()),
			zap.Duration("latency", duration),
		)
	}
}
