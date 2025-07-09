package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 连接MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("连接MongoDB失败:", err)
	}
	defer client.Disconnect(ctx)

	// 测试连接
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB连接测试失败:", err)
	}
	fmt.Println("✅ MongoDB连接成功")

	db := client.Database("auth_center")

	// 检查集合
	collections := []string{"users", "roles", "permissions", "categories", "tags", "ai_sessions", "sessions"}
	fmt.Println("\n📋 检查数据库集合:")
	for _, collName := range collections {
		count, err := db.Collection(collName).CountDocuments(ctx, bson.M{})
		if err != nil {
			fmt.Printf("❌ %s: 检查失败 - %v\n", collName, err)
		} else {
			fmt.Printf("✅ %s: %d 条记录\n", collName, count)
		}
	}

	// 检查权限数据
	fmt.Println("\n🔐 检查权限数据:")
	cursor, err := db.Collection("permissions").Find(ctx, bson.M{}, options.Find().SetLimit(5))
	if err != nil {
		fmt.Printf("❌ 查询权限失败: %v\n", err)
	} else {
		var permissions []bson.M
		if err = cursor.All(ctx, &permissions); err != nil {
			fmt.Printf("❌ 读取权限失败: %v\n", err)
		} else {
			for _, perm := range permissions {
				fmt.Printf("  - %s: %s\n", perm["name"], perm["description"])
			}
		}
	}

	// 检查角色数据
	fmt.Println("\n👥 检查角色数据:")
	cursor, err = db.Collection("roles").Find(ctx, bson.M{}, options.Find().SetLimit(5))
	if err != nil {
		fmt.Printf("❌ 查询角色失败: %v\n", err)
	} else {
		var roles []bson.M
		if err = cursor.All(ctx, &roles); err != nil {
			fmt.Printf("❌ 读取角色失败: %v\n", err)
		} else {
			for _, role := range roles {
				fmt.Printf("  - %s: %s\n", role["name"], role["description"])
			}
		}
	}

	// 检查索引
	fmt.Println("\n📊 检查数据库索引:")
	checkIndexes := map[string]string{
		"users":       "username, email, phone",
		"roles":       "name",
		"permissions": "name, resource",
		"categories":  "name",
		"tags":        "name",
	}

	for collName, expectedIndexes := range checkIndexes {
		cursor, err := db.Collection(collName).Indexes().List(ctx)
		if err != nil {
			fmt.Printf("❌ %s: 获取索引失败 - %v\n", collName, err)
			continue
		}

		var indexes []bson.M
		if err = cursor.All(ctx, &indexes); err != nil {
			fmt.Printf("❌ %s: 读取索引失败 - %v\n", collName, err)
			continue
		}

		fmt.Printf("✅ %s: %d 个索引 (期望: %s)\n", collName, len(indexes), expectedIndexes)
		for _, index := range indexes {
			if name, ok := index["name"].(string); ok && name != "_id_" {
				fmt.Printf("    - %s\n", name)
			}
		}
	}

	fmt.Println("\n🎉 数据库验证完成!")
}
