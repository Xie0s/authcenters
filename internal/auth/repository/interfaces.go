package repository

import (
	"context"

	"authcenter/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRepository 用户数据访问接口
type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByPhone(ctx context.Context, phone string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, limit, offset int) ([]*models.User, error)
}

// SessionRepository 会话数据访问接口
type SessionRepository interface {
	Create(ctx context.Context, session *models.Session) error
	GetBySessionID(ctx context.Context, sessionID string) (*models.Session, error)
	GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.Session, error)
	Update(ctx context.Context, session *models.Session) error
	RevokeUserSessions(ctx context.Context, userID primitive.ObjectID) error
	RevokeSession(ctx context.Context, sessionID string) error
	CleanupExpiredSessions(ctx context.Context) error
}
