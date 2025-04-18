package auth

import (
	"database/sql"
	"time"

	"task-tracker/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

//! \struct Service
//! \brief Encapsulates authentication business logic.
type Service struct {
	db     *sql.DB
	secret []byte
	Logger *zap.Logger
}

//! \fn NewService(db *sql.DB, secret string, logger *zap.Logger) *Service
//! \brief Initializes a new authentication service.
//! \param db Database connection.
//! \param secret JWT secret key.
//! \param logger Logger instance.
//! \return Pointer to initialized Service.
func NewService(db *sql.DB, secret string, logger *zap.Logger) *Service {
	return &Service{
		db:     db,
		secret: []byte(secret),
		Logger: logger,
	}
}

//! \fn Register(user *models.User) (int, error)
//! \brief Creates a new user in the database.
//! \param user User data to register.
//! \return User ID and error (if any).
func (s *Service) Register(user *models.User) (int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		s.Logger.Error("Failed to hash password", zap.Error(err))
		return 0, err
	}
	query := `INSERT INTO users (username, password_hash, email) VALUES ($1, $2, $3) RETURNING id`
	var userID int
	err = s.db.QueryRow(query, user.Username, string(hash), user.Email).Scan(&userID)
	if err != nil {
		s.Logger.Error("Failed to register user", zap.Error(err))
		return 0, err
	}

	s.Logger.Info("User registered", zap.Int("user_id", userID))
	return userID, nil
}

//! \fn Login(username, password string) (string, string, error)
//! \brief Authenticates a user and generates tokens.
//! \param username User's username.
//! \param password User's password.
//! \return Access token, refresh token, and error (if any).
func (s *Service) Login(username, password string) (string, string, error) {
	var user models.User
	query := `SELECT id, username, password_hash FROM users WHERE username = $1`
	err := s.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		s.Logger.Warn("User not found", zap.Error(err))
		return "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		s.Logger.Warn("Invalid password", zap.Error(err))
		return "", "", err
	}

	accessToken, err := s.generateAccessToken(user.ID, user.Username)
	if err != nil {
		s.Logger.Error("Failed to generate access token", zap.Error(err))
		return "", "", err
	}

	refreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		s.Logger.Error("Failed to generate refresh token", zap.Error(err))
		return "", "", err
	}

	s.Logger.Info("User logged in", zap.Int("user_id", user.ID))
	return accessToken, refreshToken, nil
}

//! \fn Refresh(refreshToken string) (string, error)
//! \brief Generates a new access token using a refresh token.
//! \param refreshToken Refresh token to validate.
//! \return New access token and error (if any).
func (s *Service) Refresh(refreshToken string) (string, error) {
	var userID int
	var expiresAt time.Time
	query := `SELECT user_id, expires_at FROM refresh_tokens WHERE token = $1`
	err := s.db.QueryRow(query, refreshToken).Scan(&userID, &expiresAt)
	if err != nil || expiresAt.Before(time.Now()) {
		s.Logger.Warn("Invalid or expired refresh token", zap.Error(err))
		return "", err
	}

	var username string
	query = `SELECT username FROM users WHERE id = $1`
	err = s.db.QueryRow(query, userID).Scan(&username)
	if err != nil {
		s.Logger.Error("Failed to fetch username", zap.Error(err))
		return "", err
	}

	accessToken, err := s.generateAccessToken(userID, username)
	if err != nil {
		s.Logger.Error("Failed to generate access token", zap.Error(err))
		return "", err
	}

	s.Logger.Info("Token refreshed", zap.Int("user_id", userID))
	return accessToken, nil
}

//! \fn VerifyToken(tokenString string) (*TokenClaims, error)
//! \brief Validates a JWT token.
//! \param tokenString JWT token to verify.
//! \return Token claims and error (if any).
func (s *Service) VerifyToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil || !token.Valid {
		s.Logger.Warn("Invalid token", zap.Error(err))
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		s.Logger.Warn("Invalid token claims")
		return nil, jwt.ErrInvalidKey
	}

	return claims, nil
}

//! \fn generateAccessToken(userID int, username string) (string, error)
//! \brief Generates a JWT access token.
//! \param userID User ID.
//! \param username User's username.
//! \return Signed access token and error (if any).
func (s *Service) generateAccessToken(userID int, username string) (string, error) {
	claims := &TokenClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

//! \fn generateRefreshToken(userID int) (string, error)
//! \brief Generates and stores a refresh token.
//! \param userID User ID.
//! \return Refresh token and error (if any).
func (s *Service) generateRefreshToken(userID int) (string, error) {
	token := uuid.New().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err := s.db.Exec(query, userID, token, expiresAt)
	if err != nil {
		return "", err
	}
	return token, nil
}