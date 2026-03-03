package middleware

import (
	"time"

	"github.com/aswinsreeraj/evntx/pkg/logger"
	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()

		c.Next()

		duration := time.Since(start)

		userID := c.GetString("user_id")

		logger.Log.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Str("ip", c.ClientIP()).
			Str("user_id", userID).
			Dur("duration", duration).
			Msg("request handled")
	}
}
