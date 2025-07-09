package repository

import (
	"authcenter/internal/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AIRepository AI数据访问接口
type AIRepository interface {
	// CreateSession 创建AI会话
	CreateSession(session *models.AISession) error

	// GetSession 通过ID获取AI会话
	GetSession(sessionID string) (*models.AISession, error)

	// GetSessionsByUser 获取用户的AI会话列表
	GetSessionsByUser(userID string, page, pageSize int) ([]*models.AISession, int64, error)

	// UpdateSession 更新AI会话
	UpdateSession(sessionID string, data *models.AISession) error

	// DeleteSession 删除AI会话
	DeleteSession(sessionID string) error

	// SaveMessage 保存会话消息
	SaveMessage(message *models.AIMessage) error

	// GetMessages 获取会话消息
	GetMessages(sessionID string, page, pageSize int) ([]*models.AIMessage, int64, error)

	// GetMessagesByTimeRange 按时间范围获取消息
	GetMessagesByTimeRange(sessionID string, startTime, endTime time.Time) ([]*models.AIMessage, error)

	// DeleteMessage 删除消息
	DeleteMessage(messageID string) error

	// UpdateSessionTitle 更新会话标题
	UpdateSessionTitle(sessionID, title string) error

	// GetRecentSessions 获取用户最近的会话
	GetRecentSessions(userID string, limit int) ([]*models.AISession, error)
}

// aiRepository AI仓储实现
type aiRepository struct {
	db                *mongo.Database
	sessionCollection *mongo.Collection
	messageCollection *mongo.Collection
}

// NewAIRepository 创建AI仓储
func NewAIRepository(db *mongo.Database) AIRepository {
	return &aiRepository{
		db:                db,
		sessionCollection: db.Collection("ai_sessions"),
		messageCollection: db.Collection("ai_messages"),
	}
}

// CreateSession 创建AI会话
func (r *aiRepository) CreateSession(session *models.AISession) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 设置创建时间
	now := time.Now()
	session.CreatedAt = now
	session.UpdatedAt = now

	// 设置过期时间（7天后）
	if session.ExpiresAt.IsZero() {
		session.ExpiresAt = now.AddDate(0, 0, 7)
	}

	// 插入会话
	result, err := r.sessionCollection.InsertOne(ctx, session)
	if err != nil {
		return err
	}

	// 设置生成的ID
	session.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetSession 通过ID获取AI会话
func (r *aiRepository) GetSession(sessionID string) (*models.AISession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var session models.AISession
	err := r.sessionCollection.FindOne(ctx, bson.M{"session_id": sessionID}).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("AI session not found")
		}
		return nil, err
	}

	return &session, nil
}

// GetSessionsByUser 获取用户的AI会话列表
func (r *aiRepository) GetSessionsByUser(userID string, page, pageSize int) ([]*models.AISession, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, 0, errors.New("invalid user ID format")
	}

	// 计算跳过的文档数
	skip := (page - 1) * pageSize

	// 设置查询选项
	findOptions := options.Find()
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64(skip))
	findOptions.SetSort(bson.D{{Key: "updated_at", Value: -1}}) // 按更新时间倒序

	// 查询会话
	cursor, err := r.sessionCollection.Find(ctx, bson.M{"user_id": userObjectID}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var sessions []*models.AISession
	if err = cursor.All(ctx, &sessions); err != nil {
		return nil, 0, err
	}

	// 获取总数
	total, err := r.sessionCollection.CountDocuments(ctx, bson.M{"user_id": userObjectID})
	if err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}

// UpdateSession 更新AI会话
func (r *aiRepository) UpdateSession(sessionID string, data *models.AISession) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 设置更新时间
	data.UpdatedAt = time.Now()

	// 创建更新文档
	updateDoc := bson.M{"$set": data}

	result, err := r.sessionCollection.UpdateOne(ctx, bson.M{"session_id": sessionID}, updateDoc)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("AI session not found")
	}

	return nil
}

// DeleteSession 删除AI会话
func (r *aiRepository) DeleteSession(sessionID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 删除会话相关的所有消息
	_, err := r.messageCollection.DeleteMany(ctx, bson.M{"session_id": sessionID})
	if err != nil {
		return err
	}

	// 删除会话
	result, err := r.sessionCollection.DeleteOne(ctx, bson.M{"session_id": sessionID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("AI session not found")
	}

	return nil
}

// SaveMessage 保存会话消息
func (r *aiRepository) SaveMessage(message *models.AIMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 设置时间戳
	if message.Timestamp.IsZero() {
		message.Timestamp = time.Now()
	}

	// 插入消息
	result, err := r.messageCollection.InsertOne(ctx, message)
	if err != nil {
		return err
	}

	// 设置生成的ID
	message.ID = result.InsertedID.(primitive.ObjectID)

	// 更新会话的最后更新时间
	_, err = r.sessionCollection.UpdateOne(
		ctx,
		bson.M{"session_id": message.SessionID},
		bson.M{"$set": bson.M{"updated_at": time.Now()}},
	)

	return err
}

// GetMessages 获取会话消息
func (r *aiRepository) GetMessages(sessionID string, page, pageSize int) ([]*models.AIMessage, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 计算跳过的文档数
	skip := (page - 1) * pageSize

	// 设置查询选项
	findOptions := options.Find()
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64(skip))
	findOptions.SetSort(bson.D{{Key: "timestamp", Value: 1}}) // 按时间正序

	// 查询消息
	cursor, err := r.messageCollection.Find(ctx, bson.M{"session_id": sessionID}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var messages []*models.AIMessage
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, 0, err
	}

	// 获取总数
	total, err := r.messageCollection.CountDocuments(ctx, bson.M{"session_id": sessionID})
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// GetMessagesByTimeRange 按时间范围获取消息
func (r *aiRepository) GetMessagesByTimeRange(sessionID string, startTime, endTime time.Time) ([]*models.AIMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"session_id": sessionID,
		"timestamp": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}

	cursor, err := r.messageCollection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*models.AIMessage
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

// DeleteMessage 删除消息
func (r *aiRepository) DeleteMessage(messageID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := r.messageCollection.DeleteOne(ctx, bson.M{"message_id": messageID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("message not found")
	}

	return nil
}

// UpdateSessionTitle 更新会话标题
func (r *aiRepository) UpdateSessionTitle(sessionID, title string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title":      title,
			"updated_at": time.Now(),
		},
	}

	result, err := r.sessionCollection.UpdateOne(ctx, bson.M{"session_id": sessionID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("AI session not found")
	}

	return nil
}

// GetRecentSessions 获取用户最近的会话
func (r *aiRepository) GetRecentSessions(userID string, limit int) ([]*models.AISession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.D{{Key: "updated_at", Value: -1}}) // 按更新时间倒序

	cursor, err := r.sessionCollection.Find(ctx, bson.M{"user_id": userObjectID}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sessions []*models.AISession
	if err = cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}
