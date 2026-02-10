package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/JasperRosales/aircraft-system-be/internal/controller"
)

func SetupUserRoutes(router *gin.RouterGroup, userCtrl *controller.UserController) {
	users := router.Group("/users")
	{
		users.POST("/register", userCtrl.Register)
		users.POST("/login", userCtrl.Login)
		users.POST("/logout", userCtrl.Logout)
		users.GET("/:id", userCtrl.GetByID)
		users.GET("", userCtrl.GetAll)
		users.PUT("/:id", userCtrl.Update)
		users.DELETE("/:id", userCtrl.Delete)
	}
}
