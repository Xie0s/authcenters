package handler

import (
	"net/http"

	"authcenter/internal/auth/service"
	"authcenter/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RefreshTokenRequest 刷新Token请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	user, err := h.authService.Register(c, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "注册失败", err.Error())
		return
	}

	response.Success(c, user)
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	tokenData, err := h.authService.Login(c, &req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "登录失败", err.Error())
		return
	}

	response.Success(c, tokenData)
}

// RefreshToken 刷新Token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 提供更友好的错误信息
		errorMsg := "参数错误"
		if err.Error() == "EOF" {
			errorMsg = "请求体不能为空，需要提供refresh_token"
		}
		response.Error(c, http.StatusBadRequest, errorMsg, err.Error())
		return
	}

	// 验证refresh_token不为空
	if req.RefreshToken == "" {
		response.Error(c, http.StatusBadRequest, "参数错误", "refresh_token不能为空")
		return
	}

	tokenData, err := h.authService.RefreshToken(c, req.RefreshToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Token刷新失败", err.Error())
		return
	}

	response.Success(c, tokenData)
}

// VerifyToken 验证Token
func (h *AuthHandler) VerifyToken(c *gin.Context) {
	var req service.VerifyTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	result, err := h.authService.VerifyToken(c, &req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Token验证失败", err.Error())
		return
	}

	response.Success(c, result)
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		response.Error(c, http.StatusBadRequest, "缺少Token", "")
		return
	}

	// 移除 Bearer 前缀
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	err := h.authService.Logout(c, token)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "登出失败", err.Error())
		return
	}

	response.Success(c, "登出成功")
}
