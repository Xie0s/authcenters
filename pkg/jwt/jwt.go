package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Manager JWT管理器接口
type Manager interface {
	GenerateAccessToken(userID, username string, roles, permissions []string) (string, *Claims, error)
	GenerateRefreshToken(userID string) (string, *Claims, error)
	ValidateAccessToken(tokenString string) (*Claims, error)
	ValidateRefreshToken(tokenString string) (*Claims, error)
}

// Claims JWT声明
type Claims struct {
	UserID      string   `json:"user_id"`
	Username    string   `json:"username,omitempty"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	TokenType   string   `json:"token_type"`    // access, refresh
	JTI         string   `json:"jti,omitempty"` // JWT ID，用于Refresh Token
	jwt.RegisteredClaims
}

// jwtManager JWT管理器实现
type jwtManager struct {
	secretKey            []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	issuer               string
}

// NewManager 创建JWT管理器
func NewManager(secretKey string, accessTokenDuration, refreshTokenDuration time.Duration, issuer string) Manager {
	return &jwtManager{
		secretKey:            []byte(secretKey),
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
		issuer:               issuer,
	}
}

// GenerateAccessToken 生成访问令牌
func (m *jwtManager) GenerateAccessToken(userID, username string, roles, permissions []string) (string, *Claims, error) {
	now := time.Now()
	expiresAt := now.Add(m.accessTokenDuration)

	claims := &Claims{
		UserID:      userID,
		Username:    username,
		Roles:       roles,
		Permissions: permissions,
		TokenType:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", nil, err
	}

	return tokenString, claims, nil
}

// GenerateRefreshToken 生成刷新令牌
func (m *jwtManager) GenerateRefreshToken(userID string) (string, *Claims, error) {
	now := time.Now()
	expiresAt := now.Add(m.refreshTokenDuration)
	jti := uuid.New().String() // 生成唯一ID

	claims := &Claims{
		UserID:    userID,
		TokenType: "refresh",
		JTI:       jti,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", nil, err
	}

	return tokenString, claims, nil
}

// ValidateAccessToken 验证访问令牌
func (m *jwtManager) ValidateAccessToken(tokenString string) (*Claims, error) {
	return m.validateToken(tokenString, "access")
}

// ValidateRefreshToken 验证刷新令牌
func (m *jwtManager) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return m.validateToken(tokenString, "refresh")
}

// validateToken 验证令牌
func (m *jwtManager) validateToken(tokenString, expectedType string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 确保签名方法是我们期望的
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return m.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("无效的Token")
	}

	// 检查Token类型
	if claims.TokenType != expectedType {
		return nil, errors.New("Token类型不匹配")
	}

	// 检查是否过期
	if claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("Token已过期")
	}

	return claims, nil
}
