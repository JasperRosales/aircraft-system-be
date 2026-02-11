package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/JasperRosales/aircraft-system-be/internal/models"
	"github.com/JasperRosales/aircraft-system-be/internal/service"
)

type PlaneController struct {
	service *service.PlaneService
}

func NewPlaneController(svc *service.PlaneService) *PlaneController {
	return &PlaneController{service: svc}
}

func (c *PlaneController) CreatePlane(ctx *gin.Context) {
	var req models.CreatePlaneRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.service.CreatePlane(ctx.Request.Context(), &req)
	if err != nil {
		if err.Error() == service.PlaneExistsErr {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

func (c *PlaneController) GetPlane(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid plane ID"})
		return
	}

	resp, err := c.service.GetPlane(ctx.Request.Context(), id)
	if err != nil {
		if err.Error() == service.PlaneNotFoundErr {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *PlaneController) GetPlaneByTailNumber(ctx *gin.Context) {
	tailNumber := ctx.Param("tail_number")
	if tailNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "tail number is required"})
		return
	}

	resp, err := c.service.GetPlaneByTailNumber(ctx.Request.Context(), tailNumber)
	if err != nil {
		if err.Error() == service.PlaneNotFoundErr {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *PlaneController) GetAllPlanes(ctx *gin.Context) {
	planes, err := c.service.GetAllPlanes(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if planes == nil {
		planes = []models.PlaneResponse{}
	}

	ctx.JSON(http.StatusOK, planes)
}

func (c *PlaneController) UpdatePlane(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid plane ID"})
		return
	}

	var req models.UpdatePlaneRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.service.UpdatePlane(ctx.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == service.PlaneNotFoundErr || err.Error() == service.PlaneExistsErr {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *PlaneController) DeletePlane(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid plane ID"})
		return
	}

	err = c.service.DeletePlane(ctx.Request.Context(), id)
	if err != nil {
		if err.Error() == service.PlaneNotFoundErr {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *PlaneController) GetPlaneWithParts(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid plane ID"})
		return
	}

	plane, parts, err := c.service.GetPlaneWithParts(ctx.Request.Context(), id)
	if err != nil {
		if err.Error() == service.PlaneNotFoundErr {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if parts == nil {
		parts = []models.PlanePartResponse{}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"plane": plane,
		"parts": parts,
	})
}
