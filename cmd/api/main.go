package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/JasperRosales/aircraft-system-be/internal/middleware"
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

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}
	router.Run(":" + port)
}
