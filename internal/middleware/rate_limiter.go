package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)
type visitor struct {
	count int 
	windowStart time.Time
}

type RateLimiter struct {
	mu sync.Mutex
	visitors map[string]*visitor
	limit    int
	window   time.Duration
}


func NewRateLimiter ( limit int , window time.Duration) *RateLimiter{
	return &RateLimiter{
		visitors : make(map[string]*visitor),
		limit: limit,
		window: window,
	}
}
func ( rl *RateLimiter) Allow(key string) bool{
	rl.mu.Lock()
	defer rl.mu.Unlock()
now := time.Now()
	v, exists := rl.visitors[key]


	if !exists{
		rl.visitors[key] = &visitor{
			count:1,
			windowStart: now,

		}

		return true
	}


	if now.Sub(v.windowStart) >= rl.window{
		v.count = 1
		v.windowStart = now

		return true
	}

	if v.count >= rl.limit {
		return false
}
v.count++ 
return true
}
func ( rl *RateLimiter) Middleware() gin.HandlerFunc{
	return func( c *gin.Context){
		key := c.ClientIP()
		if !rl.Allow(key){
c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}