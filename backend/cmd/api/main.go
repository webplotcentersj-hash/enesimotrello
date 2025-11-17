package main

import (
	"log"
	"task-board/internal/handler"
	"task-board/internal/middleware"
	"task-board/internal/repository"
	"task-board/internal/service"
	"task-board/internal/websocket"
	"task-board/pkg/config"
	"task-board/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize Redis
	_, err = database.InitializeRedis(cfg)
	if err != nil {
		log.Fatal("Failed to initialize Redis:", err)
	}

	// Initialize repositories
	boardRepo := repository.NewBoardRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	plotAIRepo := repository.NewPlotAIRepository(db)

	// Initialize services
	boardService := service.NewBoardService(boardRepo)
	taskService := service.NewTaskService(taskRepo)
	geminiService := service.NewGeminiService()
	plotAIService := service.NewPlotAIService(plotAIRepo, geminiService)
	
	// Set board repository in task service
	if taskSvc, ok := taskService.(interface{ SetBoardRepo(repository.BoardRepository) }); ok {
		taskSvc.SetBoardRepo(boardRepo)
	}

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize user repository and service
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, cfg)

	// Initialize handlers
	boardHandler := handler.NewBoardHandler(boardService)
	taskHandler := handler.NewTaskHandler(taskService)
	plotAIHandler := handler.NewPlotAIHandler(plotAIService)
	wsHandler := handler.NewWebSocketHandler(hub)

	// Setup router
	router := gin.Default()

	// CORS middleware
	router.Use(middleware.CORS())

	// Health check endpoint (no auth required)
	router.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "TaskBoard API is running",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Authentication routes (no middleware required)
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
			auth.GET("/profile", middleware.AuthMiddleware(cfg.JWTSecret), userHandler.GetProfile)
			auth.PUT("/profile", middleware.AuthMiddleware(cfg.JWTSecret), userHandler.UpdateProfile)
		}

		// Board and Task routes - support both authenticated and anonymous users
		// Try JWT auth first, fallback to anonymous
		api.Use(func(c *gin.Context) {
			// Check if Authorization header exists
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				// Try JWT authentication
				tokenString := authHeader[7:]
				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					return []byte(cfg.JWTSecret), nil
				})
				
				if err == nil && token.Valid {
					// JWT is valid, use it
					claims, ok := token.Claims.(jwt.MapClaims)
					if ok {
						if userID, ok := claims["user_id"].(float64); ok {
							c.Set("user_id", uint(userID))
							c.Next()
							return
						}
					}
				}
				// JWT invalid, fall through to anonymous
			}
			// Use anonymous authentication
			middleware.AnonymousUserMiddleware(db)(c)
		})

		// Board routes
		boards := api.Group("/boards")
		{
			boards.GET("", boardHandler.GetBoards)
			boards.POST("", boardHandler.CreateBoard)
			boards.GET("/:id", boardHandler.GetBoard)
			boards.PUT("/:id", boardHandler.UpdateBoard)
			boards.DELETE("/:id", boardHandler.DeleteBoard)
		}

		// Task routes
		tasks := api.Group("/tasks")
		{
			tasks.GET("/board/:boardId", taskHandler.GetTasks)
			tasks.POST("/board/:boardId", taskHandler.CreateTask)
			tasks.GET("/:id", taskHandler.GetTask)
			tasks.PUT("/:id", taskHandler.UpdateTask)
			tasks.DELETE("/:id", taskHandler.DeleteTask)
		}

		// WebSocket route
		api.GET("/ws", wsHandler.HandleWebSocket)

		// Plot AI routes (require authentication)
		plotAI := api.Group("/plot-ai")
		plotAI.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			plotAI.POST("/chat", plotAIHandler.SendMessage)
			plotAI.GET("/history", plotAIHandler.GetHistory)
		}
	}

	// Start server
	log.Printf("Server starting on %s:%s", cfg.Host, cfg.Port)
	if err := router.Run(cfg.Host + ":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
