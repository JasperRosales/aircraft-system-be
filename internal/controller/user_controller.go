package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/JasperRosales/aircraft-system-be/internal/models"
	"github.com/JasperRosales/aircraft-system-be/internal/service"
)

type UserController struct {
	service    *service.UserService
	jwtService *service.JWTService
}

func NewUserController(svc *service.UserService, jwtSvc *service.JWTService) *UserController {
	return &UserController{service: svc, jwtService: jwtSvc}
}

func (c *UserController) Register(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.service.Register(ctx.Request.Context(), &req)
	if err != nil {
		if err.Error() == service.UserExistsErr {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

func (c *UserController) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.service.Login(ctx.Request.Context(), &req)
	if err != nil {
		if err.Error() == service.UserNotFoundErr || err.Error() == service.InvalidPasswordErr {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.SetCookie(
		service.CookieName,
		resp.Token,
		int(c.jwtService.GetExpiryDuration().Seconds()),
		"/",
		"",
		false,
		true, // HttpOnly
	)

	ctx.JSON(http.StatusOK, gin.H{
		"user": resp.User,
	})
}

func (c *UserController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	resp, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if err.Error() == service.UserNotFoundErr {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *UserController) GetAll(ctx *gin.Context) {
	users, err := c.service.GetAll(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if users == nil {
		users = []models.UserResponse{}
	}

	ctx.JSON(http.StatusOK, users)
}

func (c *UserController) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req models.UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.service.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == service.UserNotFoundErr {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *UserController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	err = c.service.Delete(ctx.Request.Context(), id)
	if err != nil {
		if err.Error() == service.UserNotFoundErr {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *UserController) Logout(ctx *gin.Context) {
	ctx.SetCookie(
		service.CookieName,
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
