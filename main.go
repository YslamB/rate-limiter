package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter holds the rate limiters for different device IDs.
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
}

// NewRateLimiter creates a new RateLimiter instance.
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
	}
}

// GetLimiter returns a rate limiter for the given device ID.
func (rl *RateLimiter) GetLimiter(deviceID string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[deviceID]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		// Create a new rate limiter with 1 request per second and a burst of 3 requests.
		limiter = rate.NewLimiter(1, 3)
		rl.limiters[deviceID] = limiter
		rl.mu.Unlock()
	}

	return limiter
}

// RateLimiterMiddleware creates a middleware that applies rate limiting based on device ID.
func RateLimiterMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.GetHeader("X-Header-Device-Id")
		if deviceID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Missing X-Header-Device-Id"})
			c.Abort()
			return
		}

		limiter := rl.GetLimiter(deviceID)

		// Check if the rate limiter allows the request
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"message": "Rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.New()

	rl := NewRateLimiter()

	// Apply the rate limiter middleware to all routes.
	r.Use(RateLimiterMiddleware(rl))

	r.GET("/json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	r.Run() // Listen and serve on 0.0.0.0:8080
}
