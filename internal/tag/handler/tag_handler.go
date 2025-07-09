package handler

import (
	"net/http"

	"authcenter/internal/tag/service"

	"github.com/gin-gonic/gin"
)

// TagHandler 标签处理器
type TagHandler struct {
	tagService service.TagService
}

// NewTagHandler 创建标签处理器
func NewTagHandler(tagService service.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// GetTags 获取标签列表
func (h *TagHandler) GetTags(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get tags - not implemented"})
}

// GetTag 获取标签详情
func (h *TagHandler) GetTag(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get tag - not implemented"})
}

// CreateTag 创建标签
func (h *TagHandler) CreateTag(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "create tag - not implemented"})
}

// UpdateTag 更新标签
func (h *TagHandler) UpdateTag(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "update tag - not implemented"})
}

// DeleteTag 删除标签
func (h *TagHandler) DeleteTag(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "delete tag - not implemented"})
}

// SearchTags 搜索标签
func (h *TagHandler) SearchTags(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "search tags - not implemented"})
}

// GetTagDocuments 获取标签下的文档
func (h *TagHandler) GetTagDocuments(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get tag documents - not implemented"})
}

// GetPopularTags 获取热门标签
func (h *TagHandler) GetPopularTags(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get popular tags - not implemented"})
}
