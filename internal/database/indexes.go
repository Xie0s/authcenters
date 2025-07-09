package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// createIndexes 创建数据库索引
func createIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 用户集合索引
	if err := createUserIndexes(ctx); err != nil {
		return err
	}

	// 角色集合索引
	if err := createRoleIndexes(ctx); err != nil {
		return err
	}

	// 权限集合索引
	if err := createPermissionIndexes(ctx); err != nil {
		return err
	}

	// 分类集合索引
	if err := createCategoryIndexes(ctx); err != nil {
		return err
	}

	// 标签集合索引
	if err := createTagIndexes(ctx); err != nil {
		return err
	}

	// 知识库文档集合索引
	if err := createKnowledgeDocumentIndexes(ctx); err != nil {
		return err
	}

	// 会话集合索引
	if err := createSessionIndexes(ctx); err != nil {
		return err
	}

	// AI助手会话集合索引
	if err := createAISessionIndexes(ctx); err != nil {
		return err
	}

	// AI助手消息集合索引
	if err := createAIMessageIndexes(ctx); err != nil {
		return err
	}

	return nil
}

// createUserIndexes 创建用户集合索引
func createUserIndexes(ctx context.Context) error {
	collection := GetCollection("users")

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true),
		},
		{
			Keys:    bson.D{{Key: "phone", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true),
		},
		{
			Keys: bson.D{{Key: "roles.role_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// createRoleIndexes 创建角色集合索引
func createRoleIndexes(ctx context.Context) error {
	collection := GetCollection("roles")

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "level", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "permissions.name", Value: 1}},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// createPermissionIndexes 创建权限集合索引
func createPermissionIndexes(ctx context.Context) error {
	collection := GetCollection("permissions")

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "resource", Value: 1},
				{Key: "action", Value: 1},
			},
		},
		{
			Keys: bson.D{{Key: "category", Value: 1}},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// createCategoryIndexes 创建分类集合索引
func createCategoryIndexes(ctx context.Context) error {
	collection := GetCollection("categories")

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "parent_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "path", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "level", Value: 1},
				{Key: "sort_order", Value: 1},
			},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// createTagIndexes 创建标签集合索引
func createTagIndexes(ctx context.Context) error {
	collection := GetCollection("tags")

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "created_by", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "usage_count", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "last_used_at", Value: -1}},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// createKnowledgeDocumentIndexes 创建知识库文档集合索引
func createKnowledgeDocumentIndexes(ctx context.Context) error {
	collection := GetCollection("knowledge_documents")

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "title", Value: "text"},
				{Key: "content", Value: "text"},
				{Key: "summary", Value: "text"},
			},
		},
		{
			Keys: bson.D{{Key: "category.id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "tags.id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "author.id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "stats.view_count", Value: -1}},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// createSessionIndexes 创建会话集合索引
func createSessionIndexes(ctx context.Context) error {
	collection := GetCollection("sessions")

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "session_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys:    bson.D{{Key: "expires_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
		{
			Keys: bson.D{{Key: "is_revoked", Value: 1}},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// createAISessionIndexes 创建AI助手会话集合索引
func createAISessionIndexes(ctx context.Context) error {
	collection := GetCollection("ai_sessions")

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "session_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys:    bson.D{{Key: "expires_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
		{
			Keys: bson.D{{Key: "updated_at", Value: -1}},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// createAIMessageIndexes 创建AI助手消息集合索引
func createAIMessageIndexes(ctx context.Context) error {
	collection := GetCollection("ai_messages")

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "session_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "timestamp", Value: 1}},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}
