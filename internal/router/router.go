package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	aiHandler "authcenter/internal/ai/handler"
	aiRepo "authcenter/internal/ai/repository"
	aiService "authcenter/internal/ai/service"
	"authcenter/internal/auth/handler"
	authRepo "authcenter/internal/auth/repository"
	authService "authcenter/internal/auth/service"
	categoryHandler "authcenter/internal/category/handler"
	categoryRepo "authcenter/internal/category/repository"
	categoryService "authcenter/internal/category/service"
	"authcenter/internal/config"
	"authcenter/internal/middleware"
	permissionHandler "authcenter/internal/permission/handler"
	roleHandler "authcenter/internal/role/handler"
	roleRepo "authcenter/internal/role/repository"
	roleService "authcenter/internal/role/service"
	tagHandler "authcenter/internal/tag/handler"
	tagRepo "authcenter/internal/tag/repository"
	tagService "authcenter/internal/tag/service"
	userHandler "authcenter/internal/user/handler"
	userRepo "authcenter/internal/user/repository"
	userService "authcenter/internal/user/service"
	"authcenter/pkg/jwt"
)

// Setup 设置路由
func Setup(db *mongo.Database, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// 添加安全中间件
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.AuditMiddleware())
	r.Use(middleware.SecurityEventMiddleware())

	// 创建限流器
	rateLimiter := middleware.NewRateLimiter(100, time.Minute) // 每分钟100次请求
	r.Use(rateLimiter.RateLimit())

	// 创建JWT管理器
	jwtManager := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.AccessTokenExpire, cfg.JWT.RefreshTokenExpire, cfg.JWT.Issuer)

	// 创建中间件
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)
	loginRateLimiter := middleware.NewRateLimiter(50, 1*time.Minute) // 登录限流：1分钟50次（开发调试用）

	// 创建Repository
	userRepository := userRepo.NewUserRepository(db)
	sessionRepository := authRepo.NewSessionRepository(db)
	roleRepository := roleRepo.NewRoleRepository(db)
	categoryRepository := categoryRepo.NewCategoryRepository(db)
	tagRepository := tagRepo.NewTagRepository(db)
	aiRepository := aiRepo.NewAIRepository(db)

	// 创建Service
	authSvc := authService.NewAuthService(userRepository, sessionRepository, roleRepository, jwtManager)
	userSvc := userService.NewUserService(userRepository, roleRepository)
	roleSvc := roleService.NewRoleService(roleRepository)
	categorySvc := categoryService.NewCategoryService(categoryRepository)
	tagSvc := tagService.NewTagService(tagRepository)
	aiSvc := aiService.NewAIService(aiRepository)

	// 创建Handler
	authHdl := handler.NewAuthHandler(authSvc)
	userHdl := userHandler.NewUserHandler(userSvc)
	roleHdl := roleHandler.NewRoleHandler(roleSvc)
	permissionHdl := permissionHandler.NewPermissionHandler()
	categoryHdl := categoryHandler.NewCategoryHandler(categorySvc)
	tagHdl := tagHandler.NewTagHandler(tagSvc)
	aiHdl := aiHandler.NewAIHandler(aiSvc)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "AuthCenter"})
	})

	// API路由组
	api := r.Group("/api/v1")

	// 认证相关路由（无需认证）
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHdl.Register)
		auth.POST("/login", loginRateLimiter.RateLimit(), authHdl.Login) // 登录限流
		auth.POST("/refresh", authHdl.RefreshToken)
		auth.POST("/verify", authHdl.VerifyToken)
		auth.POST("/logout", authHdl.Logout)
	}

	// 需要认证的路由组
	protected := api.Group("")
	protected.Use(authMiddleware.RequireAuth())
	{
		// 用户管理
		users := protected.Group("/users")
		{
			users.GET("", authMiddleware.RequirePermission("user", "READ"), userHdl.GetUsers)
			users.GET("/:id", userHdl.GetUser)
			users.PUT("/:id", userHdl.UpdateUser)
			users.DELETE("/:id", authMiddleware.RequirePermission("user", "DELETE"), userHdl.DeleteUser)
			users.POST("/:id/roles", authMiddleware.RequirePermission("user", "MANAGE"), userHdl.AssignRole)
			users.DELETE("/:id/roles/:role_id", authMiddleware.RequirePermission("user", "MANAGE"), userHdl.RemoveRole)
			users.GET("/:id/permissions", userHdl.GetUserPermissions)
		}

		// 角色管理
		roles := protected.Group("/roles")
		roles.Use(authMiddleware.RequirePermission("role", "MANAGE"))
		{
			roles.GET("", roleHdl.GetRoles)
			roles.POST("", roleHdl.CreateRole)
			roles.GET("/:id", roleHdl.GetRole)
			roles.PUT("/:id", roleHdl.UpdateRole)
			roles.DELETE("/:id", roleHdl.DeleteRole)
			roles.POST("/:id/permissions", roleHdl.AssignPermission)
			roles.DELETE("/:id/permissions/:permission_id", roleHdl.RemovePermission)
		}

		// 权限管理
		permissions := protected.Group("/permissions")
		permissions.Use(authMiddleware.RequirePermission("permission", "MANAGE"))
		{
			permissions.GET("", permissionHdl.GetPermissions)
			permissions.POST("", permissionHdl.CreatePermission)
			permissions.GET("/:id", permissionHdl.GetPermission)
			permissions.PUT("/:id", permissionHdl.UpdatePermission)
			permissions.DELETE("/:id", permissionHdl.DeletePermission)
		}

		// 分类管理
		categories := protected.Group("/categories")
		{
			categories.GET("", categoryHdl.GetCategories)
			categories.POST("", authMiddleware.RequirePermission("category", "MANAGE"), categoryHdl.CreateCategory)
			categories.GET("/:id", categoryHdl.GetCategory)
			categories.PUT("/:id", authMiddleware.RequirePermission("category", "MANAGE"), categoryHdl.UpdateCategory)
			categories.DELETE("/:id", authMiddleware.RequirePermission("category", "MANAGE"), categoryHdl.DeleteCategory)
			categories.GET("/:id/documents", categoryHdl.GetCategoryDocuments)
		}

		// 标签管理
		tags := protected.Group("/tags")
		{
			tags.GET("", tagHdl.GetTags)
			tags.POST("", authMiddleware.RequirePermission("tag", "CREATE"), tagHdl.CreateTag)
			tags.GET("/:id", tagHdl.GetTag)
			tags.PUT("/:id", tagHdl.UpdateTag)    // 权限检查在handler内部处理
			tags.DELETE("/:id", tagHdl.DeleteTag) // 权限检查在handler内部处理
			tags.GET("/:id/documents", tagHdl.GetTagDocuments)
			tags.GET("/popular", tagHdl.GetPopularTags)
		}

		// AI助手
		ai := protected.Group("/ai")
		ai.Use(authMiddleware.RequirePermission("ai", "USE"))
		{
			ai.POST("/chat", aiHdl.Chat)
			ai.GET("/sessions", aiHdl.GetSessions)
			ai.GET("/sessions/:session_id", aiHdl.GetSession)
			ai.DELETE("/sessions/:session_id", aiHdl.DeleteSession)
		}
	}

	// 静态文件服务 - 用于测试页面
	r.Static("/test", "./test")

	return r
}
