package main

import (
	"log"

	"authcenter/internal/config"
	"authcenter/internal/database"
	"authcenter/internal/router"
	"authcenter/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 打印配置信息用于调试
	log.Printf("MongoDB URI: %s", cfg.MongoDB.URI)
	log.Printf("MongoDB Database: %s", cfg.MongoDB.Database)

	// 初始化日志
	logger.Init(cfg.Server.Mode)

	// 连接数据库
	db, err := database.Connect(cfg.MongoDB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Disconnect()

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化路由
	r := router.Setup(db, cfg)

	// 启动服务器
	logger.Info("Starting server on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
