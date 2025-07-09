package handler

import (
	"net/http"

	"authcenter/internal/category/service"

	"github.com/gin-gonic/gin"
)

// CategoryHandler 类别处理器
type CategoryHandler struct {
	categoryService service.CategoryService
}

// NewCategoryHandler 创建类别处理器
func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// GetCategories 获取类别列表
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get categories - not implemented"})
}

// GetCategory 获取类别详情
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get category - not implemented"})
}

// CreateCategory 创建类别
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "create category - not implemented"})
}

// UpdateCategory 更新类别
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "update category - not implemented"})
}

// DeleteCategory 删除类别
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "delete category - not implemented"})
}

// GetCategoryTree 获取类别树
func (h *CategoryHandler) GetCategoryTree(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get category tree - not implemented"})
}

// GetCategoryDocuments 获取类别下的文档
func (h *CategoryHandler) GetCategoryDocuments(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get category documents - not implemented"})
}
