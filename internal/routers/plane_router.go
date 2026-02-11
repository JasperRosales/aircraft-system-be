package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/JasperRosales/aircraft-system-be/internal/controller"
	"github.com/JasperRosales/aircraft-system-be/internal/middleware"
	"github.com/JasperRosales/aircraft-system-be/internal/service"
	"github.com/JasperRosales/aircraft-system-be/internal/util"
)

func SetupPlaneRoutes(router *gin.RouterGroup, planeCtrl *controller.PlaneController, planePartCtrl *controller.PlanePartController, jwtSvc *service.JWTService, logger *util.Logger) {
	// Protected routes (authentication required)
	planes := router.Group("/planes")
	planes.Use(middleware.AuthMiddleware(logger, jwtSvc))
	{
		// Plane CRUD
		planes.POST("", planeCtrl.CreatePlane)
		planes.GET("", planeCtrl.GetAllPlanes)
		planes.GET("/:id", planeCtrl.GetPlane)
		planes.GET("/tail/:tail_number", planeCtrl.GetPlaneByTailNumber)
		planes.PUT("/:id", planeCtrl.UpdatePlane)
		planes.DELETE("/:id", planeCtrl.DeletePlane)
		planes.GET("/:id/with-parts", planeCtrl.GetPlaneWithParts)

		// Plane Parts
		planes.POST("/:planeId/parts", planePartCtrl.AddPart)
		planes.GET("/:planeId/parts", planePartCtrl.GetPartsByPlane)
		planes.GET("/parts", planePartCtrl.GetAllParts)
		planes.GET("/parts/:partId", planePartCtrl.GetPart)
		planes.PUT("/parts/:partId", planePartCtrl.UpdatePart)
		planes.PUT("/parts/:partId/usage", planePartCtrl.UpdatePartUsage)
		planes.DELETE("/parts/:partId", planePartCtrl.DeletePart)

		// Maintenance Monitoring
		planes.GET("/maintenance/alerts", planePartCtrl.GetPartsNeedingMaintenance)
	}
}
