package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/JasperRosales/aircraft-system-be/internal/controller"
	"github.com/JasperRosales/aircraft-system-be/internal/middleware"
	"github.com/JasperRosales/aircraft-system-be/internal/repository"
	"github.com/JasperRosales/aircraft-system-be/internal/routers"
	"github.com/JasperRosales/aircraft-system-be/internal/service"
)

func main() {
	godotenv.Load()

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"logs/api.log"}
	logger, err := config.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	db, err := initDatabase()
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	logger.Info("Database connected successfully")

	userRepo := repository.NewUserRepository(db)
	jwtSvc := service.NewJWTService()
	userSvc := service.NewUserService(userRepo, jwtSvc)
	userCtrl := controller.NewUserController(userSvc, jwtSvc)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware(logger))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":     "Aircraft API is running...",
			"version":     "1.0.0",
			"description": "This API provides information about aircrafts, including their specifications, performance, and history.",
		})
	})

	api := router.Group("/api")
	routers.SetupUserRoutes(api, userCtrl)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	logger.Info("Starting server", zap.String("port", port))
	router.Run(":" + port)
}

func initDatabase() (*gorm.DB, error) {
	dsn := os.Getenv("GOOSE_DBSTRING")
	if dsn == "" {
		return nil, nil
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
