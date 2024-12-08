package handlers

import (
	"database/sql"
	"net/http"
	"time"
	"log"
	"wira-dashboard/models"
	"wira-dashboard/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db *sql.DB
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username exists
	var exists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", req.Username).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Check if email exists
	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	// Create user
	var userID int
	err = h.db.QueryRow(`
		INSERT INTO users (username, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id`,
		req.Username, req.Email, hashedPassword).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	// Log registration
	h.LogUserActivity(userID, "register", "New user registration", c)

	// Generate tokens
	token, err := utils.GenerateJWT(userID, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating refresh token"})
		return
	}

	// Store refresh token
	_, err = h.db.Exec(`
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)`,
		userID, refreshToken, time.Now().Add(time.Hour*24*30))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error storing refresh token"})
		return
	}

	c.JSON(http.StatusCreated, models.TokenResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		ExpiresIn:    3600,
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Login request for user: %s, has 2FA code: %v", req.Username, req.TOTPCode != "")

	// Get user from database
	var user models.User
	err := h.db.QueryRow(`
		SELECT id, username, password_hash, two_factor_enabled, two_factor_secret
		FROM users WHERE username = $1`,
		req.Username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.TwoFactorEnabled, &user.TwoFactorSecret)
	if err != nil {
		if err == sql.ErrNoRows {
			// Log failed login attempt
			h.LogUserActivity(0, "login_failed", "Failed login attempt: user not found", c)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	log.Printf("Found user: %s (ID: %d), 2FA enabled: %v", user.Username, user.ID, user.TwoFactorEnabled)

	// Verify password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		// Log failed login attempt
		h.LogUserActivity(user.ID, "login_failed", "Failed login attempt: incorrect password", c)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check 2FA if enabled
	if user.TwoFactorEnabled {
		log.Printf("2FA is enabled, code provided: %v", req.TOTPCode != "")
		if req.TOTPCode == "" {
			log.Printf("No 2FA code provided, requesting 2FA")
			c.JSON(http.StatusOK, gin.H{
				"requires_2fa": true,
				"message": "2FA code required",
			})
			return
		}
		if !utils.Validate2FACode(user.TwoFactorSecret.String, req.TOTPCode) {
			log.Printf("Invalid 2FA code provided")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid 2FA code"})
			return
		}
		log.Printf("2FA code validated successfully")
	}

	// Generate tokens
	token, err := utils.GenerateJWT(user.ID, user.Username)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		log.Printf("Error generating refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating refresh token"})
		return
	}

	// Store refresh token
	_, err = h.db.Exec(`
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)`,
		user.ID, refreshToken, time.Now().Add(time.Hour*24*30))
	if err != nil {
		log.Printf("Error storing refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error storing refresh token"})
		return
	}

	// Log successful login
	h.LogUserActivity(user.ID, "login", "User logged in successfully", c)

	log.Printf("Login successful for user: %s", req.Username)

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"refresh_token": refreshToken,
		"expires_in": 3600,
	})
}

// Setup2FA initiates 2FA setup for a user
func (h *AuthHandler) Setup2FA(c *gin.Context) {
	userID := c.GetInt("user_id")

	// Generate 2FA secret
	secret, qrURL, err := utils.Generate2FASecret()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating 2FA secret"})
		return
	}

	// Store secret temporarily (it will be confirmed in Enable2FA)
	_, err = h.db.Exec("UPDATE users SET two_factor_secret = $1 WHERE id = $2", secret, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"secret": secret,
		"qr_url": qrURL,
	})
}

// Enable2FA confirms and enables 2FA for a user
func (h *AuthHandler) Enable2FA(c *gin.Context) {
	var req models.Enable2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt("user_id")

	// Get user's temporary secret
	var secret string
	err := h.db.QueryRow("SELECT two_factor_secret FROM users WHERE id = $1", userID).Scan(&secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Validate TOTP code
	if !utils.Validate2FACode(secret, req.TOTPCode) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 2FA code"})
		return
	}

	// Enable 2FA
	_, err = h.db.Exec("UPDATE users SET two_factor_enabled = true WHERE id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA enabled successfully"})
}

// Disable2FA disables 2FA for a user
func (h *AuthHandler) Disable2FA(c *gin.Context) {
	userID := c.GetInt("user_id")

	_, err := h.db.Exec(`
		UPDATE users 
		SET two_factor_enabled = false, 
			two_factor_secret = NULL 
		WHERE id = $1`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA disabled successfully"})
}

// Get2FAStatus returns the 2FA status for the authenticated user
func (h *AuthHandler) Get2FAStatus(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var twoFactorEnabled bool
	err := h.db.QueryRow("SELECT two_factor_enabled FROM users WHERE id = $1", userID).Scan(&twoFactorEnabled)
	if err != nil {
		log.Printf("Error getting 2FA status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"two_factor_enabled": twoFactorEnabled,
	})
}

// RefreshToken generates new access token using refresh token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token required"})
		return
	}

	// Verify refresh token
	var userID int
	var username string
	err := h.db.QueryRow(`
		SELECT u.id, u.username 
		FROM users u
		JOIN refresh_tokens rt ON u.id = rt.user_id
		WHERE rt.token = $1 AND rt.expires_at > NOW()`,
		refreshToken).Scan(&userID, &username)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Generate new access token
	newToken, err := utils.GenerateJWT(userID, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newToken,
		"expires_in":   3600,
	})
}

// GetProfile handles fetching the user's profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	log.Printf("GetProfile - User ID from context: %v, exists: %v", userID, exists)
	
	if !exists {
		log.Printf("GetProfile - No user ID in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Query the database for user information
	var user struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	
	log.Printf("GetProfile - Querying database for user ID: %v", userID)
	err := h.db.QueryRow("SELECT id, username, email FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("GetProfile - User not found in database: %v", userID)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		log.Printf("GetProfile - Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}

	log.Printf("GetProfile - Successfully fetched profile for user: %v", user.Username)
	c.JSON(http.StatusOK, user)
}

// ChangePassword handles password change requests
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword    string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Get current password hash from database
	var currentHash string
	err := h.db.QueryRow("SELECT password_hash FROM users WHERE id = $1", userID).Scan(&currentHash)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to verify current password"})
		return
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(req.CurrentPassword)); err != nil {
		h.LogUserActivity(userID.(int), "password_change_failed", "Failed password change attempt: incorrect current password", c)
		c.JSON(400, gin.H{"error": "Current password is incorrect"})
		return
	}

	// Hash new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to process new password"})
		return
	}

	// Update password in database
	_, err = h.db.Exec("UPDATE users SET password_hash = $1 WHERE id = $2", string(newHash), userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update password"})
		return
	}

	// Log password change
	h.LogUserActivity(userID.(int), "password_change", "Password changed successfully", c)

	c.JSON(200, gin.H{"message": "Password updated successfully"})
}

// GetUserActivities returns the user's recent activities
func (h *AuthHandler) GetUserActivities(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Query recent activities
	rows, err := h.db.Query(`
		SELECT activity_type, description, created_at 
		FROM user_activities 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT 10`, userID)
	if err != nil {
		log.Printf("Error fetching activities: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activities"})
		return
	}
	defer rows.Close()

	var activities []struct {
		Type        string    `json:"type"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"timestamp"`
	}

	for rows.Next() {
		var activity struct {
			Type        string    `json:"type"`
			Description string    `json:"description"`
			CreatedAt   time.Time `json:"timestamp"`
		}
		if err := rows.Scan(&activity.Type, &activity.Description, &activity.CreatedAt); err != nil {
			log.Printf("Error scanning activity row: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process activities"})
			return
		}
		activities = append(activities, activity)
	}

	c.JSON(http.StatusOK, gin.H{"activities": activities})
}

// LogUserActivity logs a new user activity
func (h *AuthHandler) LogUserActivity(userID int, activityType string, description string, c *gin.Context) error {
	_, err := h.db.Exec(`
		INSERT INTO user_activities (user_id, activity_type, description, ip_address)
		VALUES ($1, $2, $3, $4)`,
		userID, activityType, description, c.ClientIP())
	if err != nil {
		log.Printf("Error logging activity: %v", err)
		return err
	}
	return nil
}
