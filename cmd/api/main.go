package main

import (
	"log"

	"task-tracker/internal/auth"
	"task-tracker/internal/config"
	"task-tracker/internal/db"
	"task-tracker/internal/middleware"
	"task-tracker/internal/tasks"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// @title Task Tracker API
// @version 1.0
// @description Simple task tracker API with authentication
// @host localhost:8080
// @BasePath /
func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// Connect to database
	dbConn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer dbConn.Close()

	// Initialize services
	authService := auth.NewService(dbConn, cfg.JWTSecret, logger)
	taskService := tasks.NewService(dbConn, logger)

	// Initialize Gin
	r := gin.Default()

	// Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))
	r.Use(middleware.MetricsMiddleware())

	// Metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Public routes
	r.POST("/register", auth.RegisterHandler(authService))
	r.POST("/login", auth.LoginHandler(authService))
	r.POST("/refresh", auth.RefreshHandler(authService))

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		protected.GET("/tasks", tasks.GetTasksHandler(taskService))
		protected.POST("/tasks", tasks.CreateTaskHandler(taskService))
		protected.GET("/tasks/:id", tasks.GetTaskHandler(taskService))
		protected.PUT("/tasks/:id", tasks.UpdateTaskHandler(taskService))
		protected.DELETE("/tasks/:id", tasks.DeleteTaskHandler(taskService))
	}

	// Start server
	logger.Info("Starting server", zap.String("port", cfg.Port))
	r.Run(":" + cfg.Port)
}

// package main

// import (
// 	"log"

// 	"task-tracker/auth"
// 	"task-tracker/config"
// 	"task-tracker/db"
// 	"task-tracker/middleware"
// 	"task-tracker/tasks"

// 	"github.com/gin-contrib/cors"
// 	"github.com/gin-gonic/gin"
// 	"github.com/prometheus/client_golang/prometheus/promhttp"
// 	"go.uber.org/zap"
// )

// // @title Task Tracker API
// // @version 1.0
// // @description Simple task tracker API with authentication
// // @host yourdomain.com
// // @BasePath /
// func main() {
// 	// Initialize logger
// 	logger, err := zap.NewProduction()
// 	if err != nil {
// 		log.Fatal("Failed to initialize logger:", err)
// 	}
// 	defer logger.Sync()

// 	// Load configuration
// 	cfg, err := config.Load()
// 	if err != nil {
// 		logger.Fatal("Failed to load config", zap.Error(err))
// 	}

// 	// Connect to database
// 	dbConn, err := db.Connect(cfg.DatabaseURL)
// 	if err != nil {
// 		logger.Fatal("Failed to connect to database", zap.Error(err))
// 	}
// 	defer dbConn.Close()

// 	// Initialize services
// 	authService := auth.NewService(dbConn, cfg.JWTSecret, logger)
// 	taskService := tasks.NewService(dbConn, logger)

// 	// Initialize Gin
// 	r := gin.Default()

// 	// Configure CORS
// 	corsConfig := cors.Config{
// 		AllowOrigins:     []string{"https://yourdomain.com"}, // Замените на ваш домен
// 		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
// 		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
// 		ExposeHeaders:    []string{"Content-Length"},
// 		AllowCredentials: true,
// 		MaxAge:           12 * 60 * 60, // 12 часов
// 	}
// 	r.Use(cors.New(corsConfig))

// 	// Middleware для метрик
// 	r.Use(middleware.MetricsMiddleware())

// 	// Metrics endpoint
// 	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

// 	// Public routes
// 	r.POST("/register", auth.RegisterHandler(authService))
// 	r.POST("/login", auth.LoginHandler(authService))
// 	r.POST("/refresh", auth.RefreshHandler(authService))

// 	// Protected routes
// 	protected := r.Group("/")
// 	protected.Use(middleware.AuthMiddleware(authService))
// 	{
// 		protected.GET("/tasks", tasks.GetTasksHandler(taskService))
// 		protected.POST("/tasks", tasks.CreateTaskHandler(taskService))
// 		protected.GET("/tasks/:id", tasks.GetTaskHandler(taskService))
// 		protected.PUT("/tasks/:id", tasks.UpdateTaskHandler(taskService))
// 		protected.DELETE("/tasks/:id", tasks.DeleteTaskHandler(taskService))
// 	}

// 	// Start server with HTTPS
// 	certFile := "/etc/letsencrypt/live/yourdomain.com/fullchain.pem" // Замените на ваш путь
// 	keyFile := "/etc/letsencrypt/live/yourdomain.com/privkey.pem"   // Замените на ваш путь
// 	logger.Info("Starting server with HTTPS", zap.String("port", ":443"), zap.String("cert", certFile), zap.String("key", keyFile))
// 	if err := r.RunTLS(":443", certFile, keyFile); err != nil {
// 		logger.Fatal("Failed to start HTTPS server", zap.Error(err))
// 	}
// }