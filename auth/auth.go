package auth

import (
	"net/http"

	"task-tracker/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

//! \struct TokenClaims
//! \brief Defines the structure for JWT claims.
type TokenClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

//! \fn RegisterHandler(s *Service) gin.HandlerFunc
//! \brief Creates a Gin handler for user registration.
//! \param s Authentication service instance.
//! \return Gin handler function.
func RegisterHandler(s *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			s.Logger.Warn("Invalid request body", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		validate := validator.New()
		if err := validate.Struct(&user); err != nil {
			s.Logger.Warn("Validation failed", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, err := s.Register(&user)
		if err != nil {
			s.Logger.Error("Failed to register user", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User registered", "user_id": userID})
	}
}

//! \fn LoginHandler(s *Service) gin.HandlerFunc
//! \brief Creates a Gin handler for user login.
//! \param s Authentication service instance.
//! \return Gin handler function.
func LoginHandler(s *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Username string `json:"username" validate:"required"`
			Password string `json:"password" validate:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			s.Logger.Warn("Invalid request body", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		validate := validator.New()
		if err := validate.Struct(&input); err != nil {
			s.Logger.Warn("Validation failed", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		accessToken, refreshToken, err := s.Login(input.Username, input.Password)
		if err != nil {
			s.Logger.Warn("Invalid credentials", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
	}
}

//! \fn RefreshHandler(s *Service) gin.HandlerFunc
//! \brief Creates a Gin handler for refreshing access tokens.
//! \param s Authentication service instance.
//! \return Gin handler function.
func RefreshHandler(s *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			RefreshToken string `json:"refresh_token" validate:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			s.Logger.Warn("Invalid request body", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		validate := validator.New()
		if err := validate.Struct(&input); err != nil {
			s.Logger.Warn("Validation failed", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		accessToken, err := s.Refresh(input.RefreshToken)
		if err != nil {
			s.Logger.Warn("Invalid or expired refresh token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
	}
}