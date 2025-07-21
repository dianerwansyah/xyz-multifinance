package middleware

import (
	"net/http"
	"strings"
	"time"

	"xyz-multifinance/config"
	"xyz-multifinance/logger"
	"xyz-multifinance/pkg/jwtutil"

	"github.com/gin-gonic/gin"
)

func Auth(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: no token provided"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwtutil.ValidateToken(tokenString)
		if err != nil {
			logger.Log.Warnf("failed to validate token: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized: invalid token",
				"details": err.Error(),
			})
			return
		}

		// Check expiry
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: token expired"})
				return
			}
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok || userIDFloat == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: 'user_id' claim missing or invalid"})
			return
		}
		userID := uint(userIDFloat)
		c.Set("user_id", userID)

		role, _ := claims["role"].(string)
		c.Set("role", role)

		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleInterface, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "role claim missing"})
			return
		}
		role, ok := roleInterface.(string)
		if !ok || role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: admin only"})
			return
		}
		c.Next()
	}
}
