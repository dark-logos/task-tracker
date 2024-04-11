package middleware

import (
	"strings"

	"task-tracker/auth"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//! \fn AuthMiddleware(s *auth.Service) gin.HandlerFunc
//! \brief Verifies JWT token in request headers.
//! \param s Authentication service instance.
//! \return Gin middleware function.
func AuthMiddleware(s *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			s.Logger.Warn("Authorization header missing")
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}

		claims, err := s.VerifyToken(tokenString)
		if err != nil {
			s.Logger.Warn("Invalid token", zap.Error(err))
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}