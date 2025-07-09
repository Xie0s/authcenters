package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"authcenter/pkg/logger"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware 日志中间件
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 自定义日志格式
		logData := map[string]interface{}{
			"timestamp":    param.TimeStamp.Format(time.RFC3339),
			"status":       param.StatusCode,
			"latency":      param.Latency.String(),
			"client_ip":    param.ClientIP,
			"method":       param.Method,
			"path":         param.Path,
			"user_agent":   param.Request.UserAgent(),
			"error":        param.ErrorMessage,
		}

		// 添加请求ID（如果存在）
		if requestID := param.Request.Header.Get("X-Request-ID"); requestID != "" {
			logData["request_id"] = requestID
		}

		logJSON, _ := json.Marshal(logData)
		return string(logJSON) + "\n"
	})
}

// AuditMiddleware 审计日志中间件
func AuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()

		// 读取请求体（用于审计）
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 处理请求
		c.Next()

		// 记录审计日志
		duration := time.Since(start)
		
		auditLog := map[string]interface{}{
			"timestamp":    start.Format(time.RFC3339),
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
			"query":        c.Request.URL.RawQuery,
			"status":       c.Writer.Status(),
			"duration_ms":  duration.Milliseconds(),
			"client_ip":    c.ClientIP(),
			"user_agent":   c.Request.UserAgent(),
		}

		// 添加用户信息（如果已认证）
		if userID, exists := c.Get("user_id"); exists {
			auditLog["user_id"] = userID
		}
		if username, exists := c.Get("username"); exists {
			auditLog["username"] = username
		}

		// 添加请求ID
		if requestID, exists := c.Get("request_id"); exists {
			auditLog["request_id"] = requestID
		}

		// 对于敏感操作，记录请求体（排除密码等敏感信息）
		if shouldLogRequestBody(c.Request.Method, c.Request.URL.Path) {
			if len(requestBody) > 0 {
				var bodyMap map[string]interface{}
				if err := json.Unmarshal(requestBody, &bodyMap); err == nil {
					// 移除敏感字段
					sanitizedBody := sanitizeRequestBody(bodyMap)
					auditLog["request_body"] = sanitizedBody
				}
			}
		}

		// 记录错误信息
		if len(c.Errors) > 0 {
			auditLog["errors"] = c.Errors.String()
		}

		// 输出审计日志
		logger.Info("audit", auditLog)
	}
}

// shouldLogRequestBody 判断是否应该记录请求体
func shouldLogRequestBody(method, path string) bool {
	// 对于POST、PUT、DELETE操作记录请求体
	if method == "POST" || method == "PUT" || method == "DELETE" {
		return true
	}
	return false
}

// sanitizeRequestBody 清理请求体中的敏感信息
func sanitizeRequestBody(body map[string]interface{}) map[string]interface{} {
	sensitiveFields := []string{
		"password", "token", "secret", "key", "credential",
		"access_token", "refresh_token", "authorization",
	}

	sanitized := make(map[string]interface{})
	for k, v := range body {
		isSensitive := false
		for _, field := range sensitiveFields {
			if k == field {
				sanitized[k] = "***REDACTED***"
				isSensitive = true
				break
			}
		}
		if !isSensitive {
			sanitized[k] = v
		}
	}
	return sanitized
}

// SecurityEventMiddleware 安全事件记录中间件
func SecurityEventMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 记录安全相关事件
		status := c.Writer.Status()
		path := c.Request.URL.Path

		// 记录认证失败
		if status == 401 && (path == "/api/v1/auth/login" || path == "/api/v1/auth/verify") {
			securityEvent := map[string]interface{}{
				"event_type":   "auth_failure",
				"timestamp":    time.Now().Format(time.RFC3339),
				"client_ip":    c.ClientIP(),
				"user_agent":   c.Request.UserAgent(),
				"path":         path,
				"status":       status,
			}
			
			if requestID, exists := c.Get("request_id"); exists {
				securityEvent["request_id"] = requestID
			}

			logger.Warn("security_event", securityEvent)
		}

		// 记录权限不足
		if status == 403 {
			securityEvent := map[string]interface{}{
				"event_type":   "access_denied",
				"timestamp":    time.Now().Format(time.RFC3339),
				"client_ip":    c.ClientIP(),
				"user_agent":   c.Request.UserAgent(),
				"path":         path,
				"status":       status,
			}

			if userID, exists := c.Get("user_id"); exists {
				securityEvent["user_id"] = userID
			}
			if requestID, exists := c.Get("request_id"); exists {
				securityEvent["request_id"] = requestID
			}

			logger.Warn("security_event", securityEvent)
		}

		// 记录频率限制
		if status == 429 {
			securityEvent := map[string]interface{}{
				"event_type":   "rate_limit_exceeded",
				"timestamp":    time.Now().Format(time.RFC3339),
				"client_ip":    c.ClientIP(),
				"user_agent":   c.Request.UserAgent(),
				"path":         path,
				"status":       status,
			}

			if requestID, exists := c.Get("request_id"); exists {
				securityEvent["request_id"] = requestID
			}

			logger.Warn("security_event", securityEvent)
		}
	}
}
