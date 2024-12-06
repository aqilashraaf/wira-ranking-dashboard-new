package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	"time"
)

// Define a visitor struct to hold the rate limiter and last seen time
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Define a struct to hold visitors
type rateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	r        rate.Limit
	b        int
}

// Create a new rate limiter
var limiter = &rateLimiter{
	visitors: make(map[string]*visitor),
	r:        1,    // 1 request per second
	b:        5,    // burst of 5 requests
}

// Clean up old entries every minute
func init() {
	go func() {
		for {
			time.Sleep(time.Minute)
			limiter.mu.Lock()
			for ip, v := range limiter.visitors {
				if time.Since(v.lastSeen) > 3*time.Minute {
					delete(limiter.visitors, ip)
				}
			}
			limiter.mu.Unlock()
		}
	}()
}

// getVisitor creates a new rate limiter for a visitor if it doesn't exist
func (rl *rateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.r, rl.b)
		rl.visitors[ip] = &visitor{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

// RateLimit middleware function
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := limiter.getVisitor(ip)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
