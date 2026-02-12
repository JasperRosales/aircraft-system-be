package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	logger := util.NewLogger()
	db, err := initDatabase(logger)
	if err != nil {
		logger.Fatal("Failed to initialize database", "error", err)
	}
	logger.Info("Database connected successfully")

	userRepo := repository.NewUserRepository(db)
	planeRepo := repository.NewPlaneRepository(db)
	planePartRepo := repository.NewPlanePartRepository(db)
	jwtSvc := service.NewJWTService()
	userSvc := service.NewUserService(userRepo, jwtSvc, logger)
	planeSvc := service.NewPlaneService(planeRepo, logger)
	planePartSvc := service.NewPlanePartService(planeRepo, planePartRepo, logger)
	userCtrl := controller.NewUserController(userSvc, jwtSvc)
	planeCtrl := controller.NewPlaneController(planeSvc)
	planePartCtrl := controller.NewPlanePartController(planePartSvc)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())
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
	routers.SetupPlaneRoutes(api, planeCtrl, planePartCtrl, jwtSvc, logger)

	port := os.Getenv("PORT")
	if port == "" {
 	   	port = "8080"
    	logger.Warn("PORT not set, defaulting to 8080")
	}

	router.Run(":" + port)

}

func initDatabase(logger *util.Logger) (*gorm.DB, error) {
	dsn := os.Getenv("GOOSE_DBSTRING")
	if dsn == "" {
		logger.Warn("GOOSE_DBSTRING not set, database connection skipped")
		return nil, nil
	}

	gormConfig := &gorm.Config{
		Logger: util.NewGormLogger(),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	return db, nil
}
