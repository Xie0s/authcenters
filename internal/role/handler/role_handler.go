package handler

import (
	"net/http"

	"authcenter/internal/role/service"

	"github.com/gin-gonic/gin"
)

// RoleHandler 角色处理器
type RoleHandler struct {
	roleService service.RoleService
}

// NewRoleHandler 创建角色处理器
func NewRoleHandler(roleService service.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// GetRoles 获取角色列表
func (h *RoleHandler) GetRoles(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get roles - not implemented"})
}

// GetRole 获取角色详情
func (h *RoleHandler) GetRole(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get role - not implemented"})
}

// CreateRole 创建角色
func (h *RoleHandler) CreateRole(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "create role - not implemented"})
}

// UpdateRole 更新角色
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "update role - not implemented"})
}

// DeleteRole 删除角色
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "delete role - not implemented"})
}

// AssignPermission 为角色分配权限
func (h *RoleHandler) AssignPermission(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "assign permission - not implemented"})
}

// RemovePermission 移除角色权限
func (h *RoleHandler) RemovePermission(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "remove permission - not implemented"})
}
