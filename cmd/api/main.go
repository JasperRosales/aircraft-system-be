package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/JasperRosales/aircraft-system-be/internal/controller"
	"github.com/JasperRosales/aircraft-system-be/internal/middleware"
	"github.com/JasperRosales/aircraft-system-be/internal/repository"
	"github.com/JasperRosales/aircraft-system-be/internal/routers"
	"github.com/JasperRosales/aircraft-system-be/internal/service"
	"github.com/JasperRosales/aircraft-system-be/internal/util"
)

func main() {
	godotenv.Load()

	// Ensure log directory exists
	if err := util.EnsureLogDirectory(); err != nil {
		log.Printf("Warning: Failed to create log directory: %v", err)
	}

	logger, err := util.NewLogger()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	db, err := initDatabase(logger)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	logger.Info("Database connected successfully")

	userRepo := repository.NewUserRepository(db)
	jwtSvc := service.NewJWTService()
	userSvc := service.NewUserService(userRepo, jwtSvc, logger)
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
	routers.SetupUserRoutes(api, userCtrl, jwtSvc, logger)

	port := os.Getenv("PORT")
	if port == "" {
		logger.Fatal("PORT environment variable not set")
	}

	logger.Info("Starting server", zap.String("port", port))
	router.Run(":" + port)
}

func initDatabase(logger *util.Logger) (*gorm.DB, error) {
	dsn := os.Getenv("GOOSE_DBSTRING")
	if dsn == "" {
		logger.Warn("GOOSE_DBSTRING not set, database connection skipped")
		return nil, nil
	}

	// Create a GORM logger that writes to Zap
	gormConfig := &gorm.Config{
		Logger: util.NewGormLogger(logger.Logger),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	return db, nil
}
