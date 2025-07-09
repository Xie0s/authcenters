package middleware

import (
	"net/http"
	"sync"
	"time"

	"authcenter/pkg/response"

	"github.com/gin-gonic/gin"
)

// RateLimiter 限流器结构
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int           // 限制次数
	window   time.Duration // 时间窗口
}

// NewRateLimiter 创建限流器
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// RateLimit 限流中间件
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		rl.mutex.Lock()
		defer rl.mutex.Unlock()

		now := time.Now()

		// 获取客户端的请求记录
		requests, exists := rl.requests[clientIP]
		if !exists {
			requests = []time.Time{}
		}

		// 清理过期的请求记录
		var validRequests []time.Time
		for _, reqTime := range requests {
			if now.Sub(reqTime) < rl.window {
				validRequests = append(validRequests, reqTime)
			}
		}

		// 检查是否超过限制
		if len(validRequests) >= rl.limit {
			response.Error(c, http.StatusTooManyRequests, "请求过于频繁", "请稍后再试")
			c.Abort()
			return
		}

		// 添加当前请求
		validRequests = append(validRequests, now)
		rl.requests[clientIP] = validRequests

		c.Next()
	}
}

// LoginRateLimit 登录专用限流中间件（更严格）
func (rl *RateLimiter) LoginRateLimit() gin.HandlerFunc {
	loginLimiter := NewRateLimiter(5, 15*time.Minute) // 15分钟内最多5次登录尝试
	return loginLimiter.RateLimit()
}

// ResetClientIP 重置特定客户端IP的限流记录
func (rl *RateLimiter) ResetClientIP(clientIP string) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	delete(rl.requests, clientIP)
}

// ResetAll 重置所有限流记录
func (rl *RateLimiter) ResetAll() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	rl.requests = make(map[string][]time.Time)
}

// CleanupExpiredRequests 清理过期的请求记录（定期调用）
func (rl *RateLimiter) CleanupExpiredRequests() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	for clientIP, requests := range rl.requests {
		var validRequests []time.Time
		for _, reqTime := range requests {
			if now.Sub(reqTime) < rl.window {
				validRequests = append(validRequests, reqTime)
			}
		}

		if len(validRequests) == 0 {
			delete(rl.requests, clientIP)
		} else {
			rl.requests[clientIP] = validRequests
		}
	}
}
