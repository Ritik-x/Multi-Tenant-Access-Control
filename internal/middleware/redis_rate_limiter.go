package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// import "team-access-control/internal/redis"
type RedisRateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
	prefix string
}

func NewRedisRateLimiter( client *redis.Client  , limit int , window time.Duration ,prefix string) *RedisRateLimiter {
	return &RedisRateLimiter{
		client: client,
		limit:  limit,
		window: window,
		prefix:prefix,
	}
}

func ( r1 *RedisRateLimiter) Allow(
	ctx context.Context , key string , 
) (bool , error){


	script := redis.NewScript(`
	local count = redis.call("INCR", KEYS[1])
	if count == 1 then
			redis.call("EXPIRE", KEYS[1], ARGV[1])
			end

		return count
	`)
	result, err := script.Run(
		ctx,
		r1.client,
		[]string{key},
		int(r1.window.Seconds()),
	).Int64()
	if err != nil {
		return false, err
	}

	if result > int64(r1.limit) {
		return false, nil
	}

	return true, nil
}

func (rl *RedisRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		key := fmt.Sprintf(
			"rate_limit:%s:%s",
			rl.prefix,
		c.ClientIP(),
		)
		allowed, err := rl.Allow(
			c.Request.Context(),
			key,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "rate limiter error",
			})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}