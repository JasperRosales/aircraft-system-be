package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/JasperRosales/aircraft-system-be/internal/controller"
	"github.com/JasperRosales/aircraft-system-be/internal/middleware"
	"github.com/JasperRosales/aircraft-system-be/internal/service"
	"github.com/JasperRosales/aircraft-system-be/internal/util"
)

func SetupUserRoutes(router *gin.RouterGroup, userCtrl *controller.UserController, jwtSvc *service.JWTService, logger *util.Logger) {
	// Public routes (no authentication required)
	users := router.Group("/users")
	users.POST("/register", userCtrl.Register)
	users.POST("/login", userCtrl.Login)
	users.POST("/logout", userCtrl.Logout)

	// Protected routes (authentication required)
	protected := users.Group("")
	protected.Use(middleware.AuthMiddleware(logger, jwtSvc))
	{
		protected.GET("/:id", userCtrl.GetByID)
		protected.GET("", userCtrl.GetAll)
		protected.PUT("/:id", userCtrl.Update)
		protected.DELETE("/:id", userCtrl.Delete)
	}
}
