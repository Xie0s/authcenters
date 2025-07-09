package handler

import (
	"net/http"

	"authcenter/internal/ai/service"

	"github.com/gin-gonic/gin"
)

// AIHandler AI处理器
type AIHandler struct {
	aiService service.AIService
}

// NewAIHandler 创建AI处理器
func NewAIHandler(aiService service.AIService) *AIHandler {
	return &AIHandler{
		aiService: aiService,
	}
}

// CreateSession 创建AI会话
func (h *AIHandler) CreateSession(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "create AI session - not implemented"})
}

// GetSession 获取AI会话
func (h *AIHandler) GetSession(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get AI session - not implemented"})
}

// GetSessions 获取AI会话列表
func (h *AIHandler) GetSessions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get AI sessions - not implemented"})
}

// UpdateSession 更新AI会话
func (h *AIHandler) UpdateSession(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "update AI session - not implemented"})
}

// DeleteSession 删除AI会话
func (h *AIHandler) DeleteSession(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "delete AI session - not implemented"})
}

// Chat AI对话
func (h *AIHandler) Chat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "AI chat - not implemented"})
}
