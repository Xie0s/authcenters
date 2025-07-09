package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User 用户模型
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username     string             `bson:"username" json:"username"`
	Email        string             `bson:"email,omitempty" json:"email,omitempty"`
	Phone        string             `bson:"phone,omitempty" json:"phone,omitempty"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	Status       string             `bson:"status" json:"status"` // active, inactive, locked
	Roles        []UserRole         `bson:"roles" json:"roles"`
	Profile      UserProfile        `bson:"profile" json:"profile"`
	LoginHistory LoginHistory       `bson:"login_history" json:"login_history"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

// UserRole 用户角色
type UserRole struct {
	RoleID    primitive.ObjectID `bson:"role_id" json:"role_id"`
	RoleName  string             `bson:"role_name" json:"role_name"`
	GrantedBy primitive.ObjectID `bson:"granted_by" json:"granted_by"`
	GrantedAt time.Time          `bson:"granted_at" json:"granted_at"`
	ExpiresAt *time.Time         `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
}

// UserProfile 用户资料
type UserProfile struct {
	Avatar     string `bson:"avatar,omitempty" json:"avatar,omitempty"`
	Department string `bson:"department,omitempty" json:"department,omitempty"`
	Position   string `bson:"position,omitempty" json:"position,omitempty"`
}

// LoginHistory 登录历史
type LoginHistory struct {
	LastLoginAt time.Time `bson:"last_login_at" json:"last_login_at"`
	LoginCount  int64     `bson:"login_count" json:"login_count"`
	LastIP      string    `bson:"last_ip" json:"last_ip"`
}

// Role 角色模型
type Role struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	DisplayName string             `bson:"display_name" json:"display_name"`
	Description string             `bson:"description" json:"description"`
	Level       int                `bson:"level" json:"level"`
	Status      string             `bson:"status" json:"status"`
	Permissions []RolePermission   `bson:"permissions" json:"permissions"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// RolePermission 角色权限
type RolePermission struct {
	PermissionID primitive.ObjectID `bson:"permission_id" json:"permission_id"`
	Name         string             `bson:"name" json:"name"`
	Resource     string             `bson:"resource" json:"resource"`
	Action       string             `bson:"action" json:"action"`
}

// Permission 权限模型
type Permission struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Resource    string             `bson:"resource" json:"resource"`
	Action      string             `bson:"action" json:"action"`
	Description string             `bson:"description" json:"description"`
	Category    string             `bson:"category" json:"category"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

// Category 分类模型
type Category struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name          string               `bson:"name" json:"name"`
	ParentID      *primitive.ObjectID  `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	Path          string               `bson:"path" json:"path"`
	Level         int                  `bson:"level" json:"level"`
	SortOrder     int                  `bson:"sort_order" json:"sort_order"`
	Description   string               `bson:"description" json:"description"`
	Status        string               `bson:"status" json:"status"`
	Children      []primitive.ObjectID `bson:"children" json:"children"`
	DocumentCount int64                `bson:"document_count" json:"document_count"`
	CreatedAt     time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time            `bson:"updated_at" json:"updated_at"`
}

// Tag 标签模型
type Tag struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name"`
	Color         string             `bson:"color" json:"color"`
	Description   string             `bson:"description" json:"description"`
	UsageCount    int64              `bson:"usage_count" json:"usage_count"`
	CreatedBy     primitive.ObjectID `bson:"created_by" json:"created_by"`
	CreatedByName string             `bson:"created_by_name" json:"created_by_name"`
	RelatedTags   []string           `bson:"related_tags" json:"related_tags"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	LastUsedAt    time.Time          `bson:"last_used_at" json:"last_used_at"`
}

// Session 会话模型
type Session struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SessionID      string             `bson:"session_id" json:"session_id"`
	UserID         primitive.ObjectID `bson:"user_id" json:"user_id"`
	DeviceInfo     DeviceInfo         `bson:"device_info" json:"device_info"`
	ExpiresAt      time.Time          `bson:"expires_at" json:"expires_at"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	LastAccessedAt time.Time          `bson:"last_accessed_at" json:"last_accessed_at"`
	IsRevoked      bool               `bson:"is_revoked" json:"is_revoked"`
}

// DeviceInfo 设备信息
type DeviceInfo struct {
	UserAgent  string `bson:"user_agent" json:"user_agent"`
	IP         string `bson:"ip" json:"ip"`
	DeviceType string `bson:"device_type" json:"device_type"` // web, mobile, api
}

// AISession AI助手会话模型
type AISession struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SessionID string             `bson:"session_id" json:"session_id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Title     string             `bson:"title" json:"title"`
	Context   string             `bson:"context" json:"context"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	ExpiresAt time.Time          `bson:"expires_at" json:"expires_at"`
}

// AIMessage AI助手消息模型
type AIMessage struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MessageID string             `bson:"message_id" json:"message_id"`
	SessionID string             `bson:"session_id" json:"session_id"`
	Role      string             `bson:"role" json:"role"` // user, assistant
	Content   string             `bson:"content" json:"content"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
	Context   string             `bson:"context,omitempty" json:"context,omitempty"`
}
