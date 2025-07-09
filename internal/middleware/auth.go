package middleware

import (
	"net/http"
	"strings"

	"authcenter/pkg/jwt"
	"authcenter/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件结构
type AuthMiddleware struct {
	jwtManager jwt.Manager
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(jwtManager jwt.Manager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

// RequireAuth 要求认证的中间件
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractToken(c)
		if token == "" {
			response.Error(c, http.StatusUnauthorized, "缺少认证Token", "")
			c.Abort()
			return
		}

		claims, err := m.jwtManager.ValidateAccessToken(token)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "无效的Token", err.Error())
			c.Abort()
			return
		}

		// 将用户信息设置到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roles", claims.Roles)
		c.Set("permissions", claims.Permissions)

		c.Next()
	}
}

// RequireRole 要求特定角色的中间件
func (m *AuthMiddleware) RequireRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")
		if !exists {
			response.Error(c, http.StatusForbidden, "无角色信息", "")
			c.Abort()
			return
		}

		userRoles, ok := roles.([]string)
		if !ok {
			response.Error(c, http.StatusForbidden, "角色信息格式错误", "")
			c.Abort()
			return
		}

		// 检查用户是否具有所需角色
		hasRole := false
		for _, userRole := range userRoles {
			for _, requiredRole := range requiredRoles {
				if userRole == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			response.Error(c, http.StatusForbidden, "权限不足", "需要角色: "+strings.Join(requiredRoles, ", "))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermission 要求特定权限的中间件
func (m *AuthMiddleware) RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, exists := c.Get("permissions")
		if !exists {
			response.Error(c, http.StatusForbidden, "无权限信息", "")
			c.Abort()
			return
		}

		userPermissions, ok := permissions.([]string)
		if !ok {
			response.Error(c, http.StatusForbidden, "权限信息格式错误", "")
			c.Abort()
			return
		}

		// 构建所需权限
		requiredPermission := resource + ":" + action

		// 检查用户是否具有所需权限
		hasPermission := false
		for _, perm := range userPermissions {
			if perm == requiredPermission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			response.Error(c, http.StatusForbidden, "权限不足", "需要权限: "+requiredPermission)
			c.Abort()
			return
		}

		c.Next()
	}
}

// extractToken 从请求中提取Token
func (m *AuthMiddleware) extractToken(c *gin.Context) string {
	// 从Authorization头提取
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		return authHeader[7:] // 移除 "Bearer " 前缀
	}

	// 从查询参数提取
	token := c.Query("token")
	if token != "" {
		return token
	}

	return ""
}
