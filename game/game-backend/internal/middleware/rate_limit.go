package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter 速率限制器
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter 创建新的速率限制器
func NewRateLimiter(requestsPerSecond float64, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(requestsPerSecond),
		burst:    burst,
	}
}

// GetLimiter 获取指定IP的限制器
func (rl *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[ip] = limiter
	}

	return limiter
}

// Cleanup 清理过期的限制器
func (rl *RateLimiter) Cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// 这里可以实现更复杂的清理逻辑
	// 目前保持简单，让GC处理
}

// RateLimitMiddleware 速率限制中间件
func RateLimitMiddleware(requestsPerSecond float64, burst int) gin.HandlerFunc {
	limiter := NewRateLimiter(requestsPerSecond, burst)

	// 启动清理协程
	go func() {
		ticker := time.NewTicker(time.Minute * 5)
		defer ticker.Stop()
		for range ticker.C {
			limiter.Cleanup()
		}
	}()

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		limiter := limiter.GetLimiter(clientIP)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// APIRateLimitMiddleware API速率限制中间件
func APIRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware(10.0, 20) // 每秒10个请求，突发20个
}

// WebSocketRateLimitMiddleware WebSocket速率限制中间件
func WebSocketRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware(5.0, 10) // 每秒5个请求，突发10个
}

// LoginRateLimitMiddleware 登录速率限制中间件
func LoginRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware(1.0, 3) // 每秒1个请求，突发3个
}
