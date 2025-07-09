package service

import (
	"context"
	"errors"
	"time"

	sessionRepo "authcenter/internal/auth/repository"
	"authcenter/internal/models"
	roleRepo "authcenter/internal/role/repository"
	userRepo "authcenter/internal/user/repository"
	"authcenter/pkg/jwt"
	"authcenter/pkg/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthService 认证服务接口
type AuthService interface {
	Register(ctx context.Context, req *RegisterRequest) (*models.User, error)
	Login(ctx context.Context, req *LoginRequest) (*TokenData, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenData, error)
	VerifyToken(ctx context.Context, req *VerifyTokenRequest) (*VerifyResult, error)
	Logout(ctx context.Context, token string) error
}

// authService 认证服务实现
type authService struct {
	userRepo    userRepo.UserRepository
	sessionRepo sessionRepo.SessionRepository
	roleRepo    roleRepo.RoleRepository
	jwtManager  jwt.Manager
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username"`
	Phone    string `json:"phone,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Code     string `json:"code,omitempty"`
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username,omitempty"` // 用户名
	Phone    string `json:"phone,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Code     string `json:"code,omitempty"`
	Type     string `json:"type"`
}

// VerifyTokenRequest 验证Token请求
type VerifyTokenRequest struct {
	Token    string `json:"token"`
	Resource string `json:"resource,omitempty"`
	Action   string `json:"action,omitempty"`
}

// TokenData Token数据
type TokenData struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
	UserID       string    `json:"user_id"`
}

// VerifyResult 验证结果
type VerifyResult struct {
	Valid       bool     `json:"valid"`
	UserID      string   `json:"user_id,omitempty"`
	Username    string   `json:"username,omitempty"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	HasAccess   bool     `json:"has_access,omitempty"`
}

// NewAuthService 创建认证服务
func NewAuthService(
	userRepo userRepo.UserRepository,
	sessionRepo sessionRepo.SessionRepository,
	roleRepo roleRepo.RoleRepository,
	jwtManager jwt.Manager,
) AuthService {
	return &authService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		roleRepo:    roleRepo,
		jwtManager:  jwtManager,
	}
}

// Register 用户注册
func (s *authService) Register(ctx context.Context, req *RegisterRequest) (*models.User, error) {
	// 验证用户名是否为空
	if req.Username == "" {
		return nil, errors.New("用户名不能为空")
	}

	// 检查用户名是否已存在
	existingUser, _ := s.userRepo.GetByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if req.Email != "" {
		existingUser, _ = s.userRepo.GetByEmail(req.Email)
		if existingUser != nil {
			return nil, errors.New("邮箱已被注册")
		}
	}

	// 检查手机号是否已存在
	if req.Phone != "" {
		existingUser, _ = s.userRepo.GetByPhone(req.Phone)
		if existingUser != nil {
			return nil, errors.New("手机号已被注册")
		}
	}

	// 创建用户
	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Phone:     req.Phone,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 如果有密码，进行加密
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = hashedPassword
	}

	// 分配默认角色 (User)
	defaultRole, err := s.roleRepo.GetByName("User")
	if err != nil {
		return nil, errors.New("获取默认角色失败")
	}

	userRole := models.UserRole{
		RoleID:    defaultRole.ID,
		RoleName:  defaultRole.Name,
		GrantedBy: primitive.NilObjectID, // 系统自动分配
		GrantedAt: time.Now(),
	}
	user.Roles = []models.UserRole{userRole}

	// 保存用户
	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *authService) Login(ctx context.Context, req *LoginRequest) (*TokenData, error) {
	var user *models.User
	var err error

	// 根据登录类型获取用户
	switch req.Type {
	case "phone":
		if req.Phone == "" || req.Code == "" {
			return nil, errors.New("手机号和验证码不能为空")
		}
		// TODO: 验证短信验证码
		user, err = s.userRepo.GetByPhone(req.Phone)
	case "username":
		if req.Username == "" || req.Password == "" {
			return nil, errors.New("用户名和密码不能为空")
		}
		user, err = s.userRepo.GetByUsername(req.Username)
		if err == nil && user != nil {
			if !utils.CheckPassword(req.Password, user.PasswordHash) {
				return nil, errors.New("密码错误")
			}
		}
	case "email", "": // 空字符串时默认为邮箱登录
		if req.Email == "" || req.Password == "" {
			return nil, errors.New("邮箱和密码不能为空")
		}
		user, err = s.userRepo.GetByEmail(req.Email)
		if err == nil && user != nil {
			if !utils.CheckPassword(req.Password, user.PasswordHash) {
				return nil, errors.New("密码错误")
			}
		}
	case "auto": // 自动识别用户名或邮箱登录
		if req.Password == "" {
			return nil, errors.New("密码不能为空")
		}

		var identifier string
		if req.Username != "" {
			identifier = req.Username
		} else if req.Email != "" {
			identifier = req.Email
		} else {
			return nil, errors.New("用户名或邮箱不能为空")
		}

		// 使用新的方法同时查询用户名和邮箱
		user, err = s.userRepo.GetByUsernameOrEmail(identifier)
		if err == nil && user != nil {
			if !utils.CheckPassword(req.Password, user.PasswordHash) {
				return nil, errors.New("密码错误")
			}
		}
	default:
		return nil, errors.New("不支持的登录类型")
	}

	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if user.Status != "active" {
		return nil, errors.New("用户已被禁用")
	}

	// 生成Token
	return s.generateTokens(ctx, user)
}

// RefreshToken 刷新Token
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*TokenData, error) {
	// 验证Refresh Token
	claims, err := s.jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("无效的Refresh Token")
	}

	// 检查会话是否存在且未被吊销
	session, err := s.sessionRepo.GetBySessionID(ctx, claims.JTI)
	if err != nil || session.IsRevoked {
		return nil, errors.New("会话已失效")
	}

	// 获取用户信息
	user, err := s.userRepo.GetByID(session.UserID.Hex())
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 生成新的Token
	return s.generateTokens(ctx, user)
}

// VerifyToken 验证Token
func (s *authService) VerifyToken(ctx context.Context, req *VerifyTokenRequest) (*VerifyResult, error) {
	// 验证Access Token
	claims, err := s.jwtManager.ValidateAccessToken(req.Token)
	if err != nil {
		return &VerifyResult{Valid: false}, err
	}

	result := &VerifyResult{
		Valid:       true,
		UserID:      claims.UserID,
		Username:    claims.Username,
		Roles:       claims.Roles,
		Permissions: claims.Permissions,
	}

	// 如果指定了资源和操作，检查权限
	if req.Resource != "" && req.Action != "" {
		result.HasAccess = s.checkPermission(claims.Permissions, req.Resource, req.Action)
	}

	return result, nil
}

// Logout 用户登出
func (s *authService) Logout(ctx context.Context, token string) error {
	// 验证Token
	claims, err := s.jwtManager.ValidateAccessToken(token)
	if err != nil {
		return err
	}

	// 吊销相关的会话
	userObjID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		return err
	}
	return s.sessionRepo.RevokeUserSessions(ctx, userObjID)
}

// generateTokens 生成Token对
func (s *authService) generateTokens(ctx context.Context, user *models.User) (*TokenData, error) {
	// 提取用户角色和权限
	roles := make([]string, len(user.Roles))
	var permissions []string
	permissionSet := make(map[string]bool) // 用于去重

	for i, role := range user.Roles {
		roles[i] = role.RoleName

		// 从角色中提取权限
		rolePermissions, err := s.roleRepo.GetRolePermissions(role.RoleID.Hex())
		if err != nil {
			continue // 忽略错误，继续处理其他角色
		}

		// 添加权限到集合中（去重）
		for _, perm := range rolePermissions {
			permKey := perm.Resource + ":" + perm.Action
			if !permissionSet[permKey] {
				permissionSet[permKey] = true
				permissions = append(permissions, permKey)
			}
		}
	}

	// 生成Access Token
	accessToken, accessClaims, err := s.jwtManager.GenerateAccessToken(user.ID.Hex(), user.Username, roles, permissions)
	if err != nil {
		return nil, err
	}

	// 生成Refresh Token
	refreshToken, refreshClaims, err := s.jwtManager.GenerateRefreshToken(user.ID.Hex())
	if err != nil {
		return nil, err
	}

	// 创建会话记录
	session := &models.Session{
		SessionID:      refreshClaims.JTI,
		UserID:         user.ID,
		ExpiresAt:      refreshClaims.ExpiresAt.Time,
		CreatedAt:      time.Now(),
		LastAccessedAt: time.Now(),
		IsRevoked:      false,
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return &TokenData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    accessClaims.ExpiresAt.Unix() - time.Now().Unix(),
		TokenType:    "Bearer",
		ExpiresAt:    accessClaims.ExpiresAt.Time,
		UserID:       user.ID.Hex(),
	}, nil
}

// checkPermission 检查权限
func (s *authService) checkPermission(permissions []string, resource, action string) bool {
	requiredPermission := resource + ":" + action
	for _, perm := range permissions {
		if perm == requiredPermission {
			return true
		}
	}
	return false
}
