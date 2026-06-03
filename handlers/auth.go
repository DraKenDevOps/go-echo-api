package handlers

import (
	"net/http"

	middleware "go-echo-api/middlewares"
	"go-echo-api/models"
	"go-echo-api/utils"

	"github.com/labstack/echo/v4"

	"go.uber.org/zap"
)

// Login handles user login
func (h *Handler) Login(c echo.Context) error {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.Bind(&body); err != nil {
		h.logger.Error("Failed to bind login credentials", zap.Error(err))
		return c.JSON(200, map[string]string{"status": "error", "message": "Invalid request"})
	}

	// Get user by email
	var user models.User
	err := h.db.QueryRow(
		"SELECT user_id, email, telephone, password, level, role_action, status, op_id FROM users WHERE username = ? AND deleted_at IS NULL",
		body.Username,
	).Scan(
		&user.UserID, &user.Email, &user.Telephone, &user.Password, &user.Level, &user.RoleAction, &user.Status, &user.OpID,
	)

	if err != nil {
		h.logger.Warn("User not found or invalid credentials", zap.String("username", body.Username))
		h.logger.Error("User not found or invalid credentials", zap.Error(err))
		return c.JSON(200, map[string]string{"status": "error", "message": "Invalid credentials"})
	}

	// Check password
	if err := utils.CheckPassword(user.Password, body.Password); err != nil {
		h.logger.Warn("Invalid password for user", zap.String("username", body.Username))
		h.logger.Error("Invalid password for user", zap.Error(err))
		return c.JSON(200, map[string]string{"status": "error", "message": "Invalid credentials"})
	}

	payload := models.JwtUser{
		UserID:     user.UserID,
		Username:   body.Username,
		Telephone:  user.Telephone,
		Email:      user.Email,
		Level:      user.Level,
		RoleAction: user.RoleAction,
	}

	token, err := utils.GenerateToken(payload, h.cfg)
	if err != nil {
		h.logger.Error("Failed to generate token", zap.Error(err))
		return c.JSON(200, map[string]string{"status": "error", "message": "Could not generate token"})
	}

	h.logger.Info("User logged in successfully", zap.String("username", body.Username), zap.Int("user_id", user.UserID))
	return c.JSON(200, map[string]any{
		"token": token,
		"user":  payload,
	})
}

// Register handles user registration
func (h *Handler) Register(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		h.logger.Error("Failed to bind user data", zap.Error(err))
		return c.JSON(200, map[string]string{"status": "error", "message": "Invalid request"})
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		h.logger.Error("Failed to hash password", zap.Error(err))
		return c.JSON(200, map[string]string{"status": "error", "message": "Could not process request"})
	}

	// Insert user
	result, err := h.db.Exec(
		"INSERT INTO users (username, telephone, email, password, level, role_action, status) VALUES (?, ?, ?, ?, ?, ?, ?)",
		user.Username, user.Telephone, user.Email, string(hashedPassword), user.Level, user.RoleAction, user.Status,
	)
	if err != nil {
		h.logger.Error("Failed to create user", zap.Error(err))
		return c.JSON(200, map[string]string{"status": "error", "message": "Could not create user"})
	}

	userID, err := result.LastInsertId()
	if err != nil {
		h.logger.Error("Failed to get last insert ID", zap.Error(err))
		return c.JSON(200, map[string]string{"status": "error", "message": "Could not create user"})
	}

	h.logger.Info("User registered successfully", zap.Int64("user_id", userID))
	return c.JSON(http.StatusCreated, map[string]any{
		"status":  "success",
		"message": "User created successfully",
	})
}

// Refresh returns the currently authenticated user's profile.
func (h *Handler) Refresh(c echo.Context) error {
	claims := middleware.GetClaims(c)
	token := middleware.ExtractToken(c)
	if token == "" {
		return c.JSON(200, map[string]string{
			"status":  "error",
			"message": "Missing authorization token",
		})
	}
	if claims == nil {
		return c.JSON(200, map[string]string{"status": "error", "message": "Unauthorized"})
	}

	var user models.User
	err := h.db.QueryRow(
		"SELECT user_id, email, telephone, password, level, role_action, status FROM users WHERE username = ? AND deleted_at IS NULL",
		claims.Username,
	).Scan(
		&user.UserID, &user.Email, &user.Telephone, &user.Password, &user.Level, &user.RoleAction, &user.Status,
	)
	if err != nil {
		h.logger.Error("User not found", zap.Int("user_id", claims.UserID), zap.Error(err))
		return c.JSON(200, map[string]string{"status": "error", "message": "User not found"})
	}

	payload := models.JwtUser{
		UserID:     user.UserID,
		Username:   claims.Username,
		Telephone:  user.Telephone,
		Email:      user.Email,
		Level:      user.Level,
		RoleAction: user.RoleAction,
	}

	h.logger.Info("User logged in successfully", zap.String("username", claims.Username), zap.Int("user_id", user.UserID))
	return c.JSON(200, map[string]any{
		"token": token,
		"user":  payload,
	})
}
