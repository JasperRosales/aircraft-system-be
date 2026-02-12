package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/JasperRosales/aircraft-system-be/internal/service"
	"github.com/JasperRosales/aircraft-system-be/internal/util"
)

func AuthMiddleware(logger *util.Logger, jwtSvc *service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(service.CookieName)
		logger.Info("Auth: Checking cookie",
			"cookie_name", service.CookieName,
			"token", token,
			"error", err,
		)

		if err != nil || token == "" {
			authHeader := c.GetHeader("Authorization")
			logger.Info("Auth: Checking header",
				"header", authHeader,
			)
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if token == "" {
			logger.Warn("Auth: No token found, rejecting request")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authentication required",
			})
			return
		}

		claims, err := jwtSvc.ValidateToken(token)
		if err != nil {
			logger.Warn("Auth: Token validation failed",
				"error", err,
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_name", claims.Name)
		c.Set("user_role", claims.Role)

		logger.Info("Auth: User authenticated",
			"user_id", claims.UserID,
			"name", claims.Name,
			"role", claims.Role,
		)
		c.Next()
	}
}

func RoleMiddleware(logger *util.Logger, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			logger.Warn("Role: User not authenticated")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authentication required",
			})
			return
		}

		if role.(string) != requiredRole && role.(string) != "admin" {
			logger.Warn("Role: Insufficient permissions",
				"user_role", role.(string),
				"required_role", requiredRole,
			)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "insufficient permissions",
			})
			return
		}

		logger.Info("Role: Access granted",
			"role", role.(string),
		)
		c.Next()
	}
}

func GetUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(int64), true
}

func GetUserName(c *gin.Context) (string, bool) {
	name, exists := c.Get("user_name")
	if !exists {
		return "", false
	}
	return name.(string), true
}

func GetUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	return role.(string), true
}
