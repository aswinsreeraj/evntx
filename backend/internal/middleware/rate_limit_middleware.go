package middleware

import (
	"net/http"
	"sync"
	"time"

	apiErrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/aswinsreeraj/evntx/pkg/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var clients = make(map[string]*client)
var mu sync.Mutex

func RateLimitMiddleware(r rate.Limit, b int) gin.HandlerFunc {
	return func(c *gin.Context) {

		ip := c.ClientIP()

		mu.Lock()

		if _, exists := clients[ip]; !exists {
			clients[ip] = &client{
				limiter: rate.NewLimiter(r, b),
			}
		}

		clients[ip].lastSeen = time.Now()
		limiter := clients[ip].limiter

		mu.Unlock()

		if !limiter.Allow() {
			response.Error(
				c,
				http.StatusTooManyRequests,
				apiErrors.RateLimitExceeded,
				"Rate limit exceeded. Please try again later.",
			)
			c.Abort()
			return
		}

		c.Next()
	}
}
