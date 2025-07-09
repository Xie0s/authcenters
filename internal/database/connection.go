package database

import (
	"context"
	"time"

	"authcenter/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var database *mongo.Database

// Connect 连接MongoDB数据库
func Connect(cfg config.MongoDBConfig) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	// 设置连接选项
	clientOptions := options.Client().
		ApplyURI(cfg.URI).
		SetMaxPoolSize(uint64(cfg.MaxPoolSize)).
		SetMinPoolSize(uint64(cfg.MinPoolSize))

	// 连接数据库
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// 测试连接
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	database = client.Database(cfg.Database)

	// 创建索引
	if err := createIndexes(); err != nil {
		return nil, err
	}

	return database, nil
}

// Disconnect 断开数据库连接
func Disconnect() error {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return client.Disconnect(ctx)
	}
	return nil
}

// GetDatabase 获取数据库实例
func GetDatabase() *mongo.Database {
	return database
}

// GetCollection 获取集合
func GetCollection(name string) *mongo.Collection {
	return database.Collection(name)
}
