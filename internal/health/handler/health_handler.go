package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"authcenter/internal/config"
	"authcenter/pkg/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	db  *mongo.Database
	cfg *config.Config
}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler(db *mongo.Database, cfg *config.Config) *HealthHandler {
	return &HealthHandler{
		db:  db,
		cfg: cfg,
	}
}

// HealthCheckResponse 健康检查响应
type HealthCheckResponse struct {
	Status    string                 `json:"status"`
	Timestamp string                 `json:"timestamp"`
	Checks    map[string]CheckResult `json:"checks"`
	Summary   Summary                `json:"summary"`
}

// CheckResult 检查结果
type CheckResult struct {
	Status      string      `json:"status"`
	Message     string      `json:"message"`
	Duration    string      `json:"duration"`
	Details     interface{} `json:"details,omitempty"`
	LastChecked string      `json:"last_checked"`
}

// Summary 汇总信息
type Summary struct {
	TotalChecks  int `json:"total_checks"`
	PassedChecks int `json:"passed_checks"`
	FailedChecks int `json:"failed_checks"`
}

// BasicHealth 基础健康检查
func (h *HealthHandler) BasicHealth(c *gin.Context) {
	response.Success(c, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "AuthCenter",
		"version":   "1.0.0",
	})
}

// DetailedHealth 详细健康检查
func (h *HealthHandler) DetailedHealth(c *gin.Context) {
	startTime := time.Now()
	checks := make(map[string]CheckResult)
	
	// 检查数据库连接
	checks["database"] = h.checkDatabase()
	
	// 检查JWT配置
	checks["jwt_config"] = h.checkJWTConfig()
	
	// 检查MongoDB集合
	checks["collections"] = h.checkCollections()
	
	// 检查配置完整性
	checks["config"] = h.checkConfigIntegrity()
	
	// 计算汇总信息
	summary := h.calculateSummary(checks)
	
	// 确定整体状态
	overallStatus := "healthy"
	if summary.FailedChecks > 0 {
		overallStatus = "unhealthy"
	}
	
	healthResponse := HealthCheckResponse{
		Status:    overallStatus,
		Timestamp: time.Now().Format(time.RFC3339),
		Checks:    checks,
		Summary:   summary,
	}
	
	// 根据状态返回适当的HTTP状态码
	statusCode := http.StatusOK
	if overallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}
	
	c.JSON(statusCode, gin.H{
		"code":    statusCode,
		"message": fmt.Sprintf("Health check completed in %v", time.Since(startTime)),
		"data":    healthResponse,
	})
}

// checkDatabase 检查数据库连接
func (h *HealthHandler) checkDatabase() CheckResult {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	err := h.db.Client().Ping(ctx, nil)
	duration := time.Since(start)
	
	if err != nil {
		return CheckResult{
			Status:      "failed",
			Message:     fmt.Sprintf("Database connection failed: %v", err),
			Duration:    duration.String(),
			LastChecked: time.Now().Format(time.RFC3339),
		}
	}
	
	return CheckResult{
		Status:      "passed",
		Message:     "Database connection successful",
		Duration:    duration.String(),
		LastChecked: time.Now().Format(time.RFC3339),
		Details: map[string]interface{}{
			"database_name": h.db.Name(),
		},
	}
}

// checkJWTConfig 检查JWT配置
func (h *HealthHandler) checkJWTConfig() CheckResult {
	start := time.Now()
	duration := time.Since(start)
	
	var issues []string
	
	if h.cfg.JWT.Secret == "" || h.cfg.JWT.Secret == "your-super-secret-jwt-key-change-this-in-production" || h.cfg.JWT.Secret == "change-this-secret-in-production" {
		issues = append(issues, "JWT secret is not set or using default value")
	}
	
	if h.cfg.JWT.AccessTokenExpire == 0 {
		issues = append(issues, "Access token expiration not configured")
	}
	
	if h.cfg.JWT.RefreshTokenExpire == 0 {
		issues = append(issues, "Refresh token expiration not configured")
	}
	
	if len(issues) > 0 {
		return CheckResult{
			Status:      "failed",
			Message:     "JWT configuration issues found",
			Duration:    duration.String(),
			LastChecked: time.Now().Format(time.RFC3339),
			Details: map[string]interface{}{
				"issues": issues,
			},
		}
	}
	
	return CheckResult{
		Status:      "passed",
		Message:     "JWT configuration is valid",
		Duration:    duration.String(),
		LastChecked: time.Now().Format(time.RFC3339),
		Details: map[string]interface{}{
			"access_token_expire":  h.cfg.JWT.AccessTokenExpire.String(),
			"refresh_token_expire": h.cfg.JWT.RefreshTokenExpire.String(),
			"issuer":               h.cfg.JWT.Issuer,
		},
	}
}

// checkCollections 检查MongoDB集合
func (h *HealthHandler) checkCollections() CheckResult {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	expectedCollections := []string{"users", "roles", "permissions", "sessions"}
	collections, err := h.db.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		return CheckResult{
			Status:      "failed",
			Message:     fmt.Sprintf("Failed to list collections: %v", err),
			Duration:    time.Since(start).String(),
			LastChecked: time.Now().Format(time.RFC3339),
		}
	}
	
	collectionMap := make(map[string]bool)
	for _, col := range collections {
		collectionMap[col] = true
	}
	
	missingCollections := []string{}
	existingCollections := []string{}
	
	for _, expected := range expectedCollections {
		if collectionMap[expected] {
			existingCollections = append(existingCollections, expected)
		} else {
			missingCollections = append(missingCollections, expected)
		}
	}
	
	status := "passed"
	message := "All required collections exist"
	
	if len(missingCollections) > 0 {
		status = "failed"
		message = "Some required collections are missing"
	}
	
	return CheckResult{
		Status:      status,
		Message:     message,
		Duration:    time.Since(start).String(),
		LastChecked: time.Now().Format(time.RFC3339),
		Details: map[string]interface{}{
			"existing_collections": existingCollections,
			"missing_collections":  missingCollections,
			"total_collections":    len(collections),
		},
	}
}

// checkConfigIntegrity 检查配置完整性
func (h *HealthHandler) checkConfigIntegrity() CheckResult {
	start := time.Now()
	duration := time.Since(start)
	
	var issues []string
	
	if h.cfg.Server.Port == "" {
		issues = append(issues, "Server port not configured")
	}
	
	if h.cfg.MongoDB.URI == "" {
		issues = append(issues, "MongoDB URI not configured")
	}
	
	if h.cfg.MongoDB.Database == "" {
		issues = append(issues, "MongoDB database name not configured")
	}
	
	if h.cfg.Security.MaxLoginAttempts <= 0 {
		issues = append(issues, "Max login attempts not properly configured")
	}
	
	if h.cfg.Security.PasswordMinLength < 8 {
		issues = append(issues, "Password minimum length is too weak (should be at least 8)")
	}
	
	if len(issues) > 0 {
		return CheckResult{
			Status:      "failed",
			Message:     "Configuration integrity issues found",
			Duration:    duration.String(),
			LastChecked: time.Now().Format(time.RFC3339),
			Details: map[string]interface{}{
				"issues": issues,
			},
		}
	}
	
	return CheckResult{
		Status:      "passed",
		Message:     "Configuration integrity check passed",
		Duration:    duration.String(),
		LastChecked: time.Now().Format(time.RFC3339),
		Details: map[string]interface{}{
			"server_mode":          h.cfg.Server.Mode,
			"max_login_attempts":   h.cfg.Security.MaxLoginAttempts,
			"password_min_length":  h.cfg.Security.PasswordMinLength,
			"session_cleanup":      h.cfg.Security.SessionCleanupInterval,
		},
	}
}

// calculateSummary 计算汇总信息
func (h *HealthHandler) calculateSummary(checks map[string]CheckResult) Summary {
	summary := Summary{
		TotalChecks: len(checks),
	}
	
	for _, check := range checks {
		if check.Status == "passed" {
			summary.PassedChecks++
		} else {
			summary.FailedChecks++
		}
	}
	
	return summary
}