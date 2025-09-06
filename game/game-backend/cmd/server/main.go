package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"game-backend/config"
	"game-backend/internal/handler"
	"game-backend/internal/middleware"
	"game-backend/internal/service"
	"game-backend/internal/websocket"
	"game-backend/pkg/database"
)

func main() {
	// 加载配置
	if err := config.LoadConfig("config/config.yaml"); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库
	if err := database.InitMySQL(); err != nil {
		log.Fatalf("初始化MySQL失败: %v", err)
	}
	defer database.Close()

	// 自动迁移数据库表结构
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 初始化Redis
	if err := database.InitRedis(); err != nil {
		log.Fatalf("初始化Redis失败: %v", err)
	}
	defer database.CloseRedis()

	// 创建WebSocket中心
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// 创建服务
	authService := service.NewAuthService(database.GetDB())
	gameService := service.NewGameService(database.GetDB())

	// 创建处理器
	authHandler := handler.NewAuthHandler(authService)
	gameHandler := handler.NewGameHandler(gameService, wsHub)

	// 设置Gin模式
	if config.AppConfig.Server.IsDebug() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	router := setupRouter(authHandler, gameHandler, wsHub)

	// 启动服务器
	serverAddr := config.AppConfig.Server.GetServerAddr()
	log.Printf("服务器启动在: %s", serverAddr)

	// 优雅关闭
	go func() {
		if err := router.Run(serverAddr); err != nil {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("服务器正在关闭...")
}

// setupRouter 设置路由
func setupRouter(authHandler *handler.AuthHandler, gameHandler *handler.GameHandler, wsHub *websocket.Hub) *gin.Engine {
	router := gin.New()

	// 中间件
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.CORSMiddleware())

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "服务运行正常",
		})
	})

	// API版本组
	v1 := router.Group("/api/v1")

	// 认证相关路由
	auth := v1.Group("/auth")
	{
		auth.POST("/login", middleware.LoginRateLimitMiddleware(), authHandler.Login)
		auth.POST("/register", middleware.LoginRateLimitMiddleware(), authHandler.Register)
		auth.POST("/logout", middleware.AuthMiddleware(), authHandler.Logout)
		auth.GET("/profile", middleware.AuthMiddleware(), authHandler.GetProfile)
		auth.PUT("/profile", middleware.AuthMiddleware(), authHandler.UpdateProfile)
		auth.POST("/refresh", middleware.AuthMiddleware(), authHandler.RefreshToken)
	}

	// 游戏相关路由
	game := v1.Group("/game")
	{
		// 公开接口
		game.GET("/status", gameHandler.GetGameStatus)
		game.GET("/history", gameHandler.GetGameHistory)
		game.GET("/leaderboard", gameHandler.GetLeaderboard)

		// 需要认证的接口
		gameAuth := game.Group("", middleware.AuthMiddleware())
		{
			gameAuth.POST("/bet", middleware.APIRateLimitMiddleware(), gameHandler.PlaceBet)
			gameAuth.POST("/cashout", middleware.APIRateLimitMiddleware(), gameHandler.Cashout)
			gameAuth.GET("/bet/history", gameHandler.GetBetHistory)
			gameAuth.GET("/stats", gameHandler.GetUserStats)
		}
	}

	// WebSocket路由
	router.GET("/ws", websocket.ServeWS(wsHub))

	return router
}
