package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/JasperRosales/aircraft-system-be/internal/models"
	"github.com/JasperRosales/aircraft-system-be/internal/service"
)

type PlanePartController struct {
	service *service.PlanePartService
}

func NewPlanePartController(svc *service.PlanePartService) *PlanePartController {
	return &PlanePartController{service: svc}
}

func (c *PlanePartController) AddPart(ctx *gin.Context) {
	planeID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid plane ID"})
		return
	}

	var req models.CreatePlanePartRequest
	req.PlaneID = planeID
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.service.AddPart(ctx.Request.Context(), &req)
	if err != nil {
		if err.Error() == service.PlaneNotFoundErrPart {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == service.PlanePartExistsErr {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

func (c *PlanePartController) GetPart(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("partId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid part ID"})
		return
	}

	resp, err := c.service.GetPart(ctx.Request.Context(), id)
	if err != nil {
		if err.Error() == service.PlanePartNotFoundErr {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *PlanePartController) GetPartsByPlane(ctx *gin.Context) {
	planeID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid plane ID"})
		return
	}
	category := ctx.Query("category")

	parts, err := c.service.GetPartsByPlane(ctx.Request.Context(), planeID, &category)
	if err != nil {
		if err.Error() == service.PlaneNotFoundErrPart {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if parts == nil {
		parts = []models.PlanePartResponse{}
	}

	ctx.JSON(http.StatusOK, parts)
}

func (c *PlanePartController) GetAllParts(ctx *gin.Context) {
	parts, err := c.service.GetAllParts(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if parts == nil {
		parts = []models.PlanePartResponse{}
	}

	ctx.JSON(http.StatusOK, parts)
}

func (c *PlanePartController) UpdatePart(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("partId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid part ID"})
		return
	}

	var req models.UpdatePlanePartRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.service.UpdatePart(ctx.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == service.PlanePartNotFoundErr || err.Error() == service.PlanePartExistsErr {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *PlanePartController) UpdatePartUsage(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("partId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid part ID"})
		return
	}

	var req models.UpdatePartUsageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.service.UpdatePartUsage(ctx.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == service.PlanePartNotFoundErr {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == service.InvalidUsageHoursErr {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *PlanePartController) DeletePart(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("partId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid part ID"})
		return
	}

	err = c.service.DeletePart(ctx.Request.Context(), id)
	if err != nil {
		if err.Error() == service.PlanePartNotFoundErr {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *PlanePartController) GetPartsNeedingMaintenance(ctx *gin.Context) {
	thresholdStr := ctx.DefaultQuery("threshold", "80")
	threshold, err := strconv.ParseFloat(thresholdStr, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid threshold value"})
		return
	}

	parts, err := c.service.GetPartsNeedingMaintenance(ctx.Request.Context(), threshold)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if parts == nil {
		parts = []models.PlanePartResponse{}
	}

	ctx.JSON(http.StatusOK, parts)
}
