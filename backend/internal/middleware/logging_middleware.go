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

		status := c.Writer.Status()
		if status < 400 {
			return
		}

		duration := time.Since(start)
		userID := c.GetString("user_id")

		event := logger.Log.Warn()
		if status >= 500 {
			event = logger.Log.Error()
		}

		event.
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", status).
			Str("ip", c.ClientIP()).
			Str("user_id", userID).
			Dur("duration", duration).
			Msg("request failed")
	}
}
