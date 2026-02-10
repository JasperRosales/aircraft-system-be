package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/JasperRosales/aircraft-system-be/internal/service"
	"github.com/JasperRosales/aircraft-system-be/internal/util"
)

func AuthMiddleware(logger *util.Logger, jwtSvc *service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(service.CookieName)
		logger.Info("Auth: Checking cookie",
			zap.String("cookie_name", service.CookieName),
			zap.String("token", token),
			zap.Error(err),
		)

		// If not in cookie, try Authorization header (Bearer token)
		if err != nil || token == "" {
			authHeader := c.GetHeader("Authorization")
			logger.Info("Auth: Checking header",
				zap.String("header", authHeader),
			)
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// If no token found, reject request
		if token == "" {
			logger.Warn("Auth: No token found, rejecting request")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authentication required",
			})
			return
		}

		// Validate token
		claims, err := jwtSvc.ValidateToken(token)
		if err != nil {
			logger.Warn("Auth: Token validation failed",
				zap.Error(err),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Set user info in context for downstream handlers
		c.Set("user_id", claims.UserID)
		c.Set("user_name", claims.Name)
		c.Set("user_role", claims.Role)

		logger.Info("Auth: User authenticated",
			zap.Int64("user_id", claims.UserID),
			zap.String("name", claims.Name),
			zap.String("role", claims.Role),
		)
		c.Next()
	}
}

// RoleMiddleware checks if the authenticated user has the required role
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
				zap.String("user_role", role.(string)),
				zap.String("required_role", requiredRole),
			)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "insufficient permissions",
			})
			return
		}

		logger.Info("Role: Access granted",
			zap.String("role", role.(string)),
		)
		c.Next()
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(int64), true
}

// GetUserName extracts user name from context
func GetUserName(c *gin.Context) (string, bool) {
	name, exists := c.Get("user_name")
	if !exists {
		return "", false
	}
	return name.(string), true
}

// GetUserRole extracts user role from context
func GetUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	return role.(string), true
}
