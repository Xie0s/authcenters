package repository

import (
	"context"
	"errors"
	"time"

	"authcenter/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// sessionRepository 会话仓储实现
type sessionRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewSessionRepository 创建会话仓储
func NewSessionRepository(db *mongo.Database) SessionRepository {
	return &sessionRepository{
		db:         db,
		collection: db.Collection("sessions"),
	}
}

// Create 创建会话
func (r *sessionRepository) Create(ctx context.Context, session *models.Session) error {
	// 设置创建时间
	now := time.Now()
	session.CreatedAt = now
	session.LastAccessedAt = now
	session.IsRevoked = false

	_, err := r.collection.InsertOne(ctx, session)
	if err != nil {
		return err
	}

	return nil
}

// GetBySessionID 通过会话ID获取会话
func (r *sessionRepository) GetBySessionID(ctx context.Context, sessionID string) (*models.Session, error) {
	var session models.Session
	err := r.collection.FindOne(ctx, bson.M{"session_id": sessionID}).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("session not found")
		}
		return nil, err
	}

	return &session, nil
}

// GetByUserID 通过用户ID获取会话列表
func (r *sessionRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.Session, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"user_id":    userID,
		"is_revoked": false,
		"expires_at": bson.M{"$gt": time.Now()}, // 未过期的会话
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sessions []*models.Session
	if err = cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

// Update 更新会话
func (r *sessionRepository) Update(ctx context.Context, session *models.Session) error {
	session.LastAccessedAt = time.Now()

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"session_id": session.SessionID},
		bson.M{"$set": session},
	)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("session not found")
	}

	return nil
}

// RevokeUserSessions 撤销用户所有会话
func (r *sessionRepository) RevokeUserSessions(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{"user_id": userID},
		bson.M{"$set": bson.M{"is_revoked": true}},
	)
	return err
}

// RevokeSession 撤销指定会话
func (r *sessionRepository) RevokeSession(ctx context.Context, sessionID string) error {
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"session_id": sessionID},
		bson.M{"$set": bson.M{"is_revoked": true}},
	)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("session not found")
	}

	return nil
}

// CleanupExpiredSessions 清理过期会话
func (r *sessionRepository) CleanupExpiredSessions(ctx context.Context) error {
	// MongoDB的TTL索引会自动清理过期文档，这里主要是手动清理被撤销的会话
	_, err := r.collection.DeleteMany(ctx, bson.M{
		"$or": []bson.M{
			{"is_revoked": true},
			{"expires_at": bson.M{"$lt": time.Now()}},
		},
	})
	return err
}
