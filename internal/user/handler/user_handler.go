package handler

import (
	"net/http"

	"authcenter/internal/user/service"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUsers 获取用户列表
func (h *UserHandler) GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get users - not implemented"})
}

// GetUser 获取用户详情
func (h *UserHandler) GetUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get user - not implemented"})
}

// UpdateUser 更新用户信息
func (h *UserHandler) UpdateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "update user - not implemented"})
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "delete user - not implemented"})
}

// AssignRole 为用户分配角色
func (h *UserHandler) AssignRole(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "assign role - not implemented"})
}

// RemoveRole 移除用户角色
func (h *UserHandler) RemoveRole(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "remove role - not implemented"})
}

// GetUserPermissions 获取用户权限
func (h *UserHandler) GetUserPermissions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get user permissions - not implemented"})
}
