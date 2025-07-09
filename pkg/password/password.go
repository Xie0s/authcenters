package password

import (
	"golang.org/x/crypto/bcrypt"
)

// Manager 密码管理器接口
type Manager interface {
	HashPassword(password string) (string, error)
	CheckPassword(password, hashedPassword string) error
}

// passwordManager 密码管理器实现
type passwordManager struct {
	cost int
}

// NewManager 创建密码管理器
func NewManager(cost int) Manager {
	// 根据需求文档，cost factor >= 12
	if cost < 12 {
		cost = 12
	}
	return &passwordManager{
		cost: cost,
	}
}

// HashPassword 加密密码
func (m *passwordManager) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), m.cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword 验证密码
func (m *passwordManager) CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
