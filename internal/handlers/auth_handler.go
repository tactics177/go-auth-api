package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/tactics177/go-auth-api/internal/models"
	"github.com/tactics177/go-auth-api/internal/repositories"
	"github.com/tactics177/go-auth-api/internal/services"
	"github.com/tactics177/go-auth-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
	"time"
)

// Register Handler
func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)

	if user.Name == "" || user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	err := services.RegisterUser(&user)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login Handler
func Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	email := strings.TrimSpace(loginData.Email)
	password := strings.TrimSpace(loginData.Password)

	if email == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	if !utils.ValidateEmail(email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	accessToken, refreshToken, err := services.AuthenticateUser(email, password)
	if err != nil {
		if err.Error() == "failed to generate token" || err.Error() == "failed to generate refresh token" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Login successful",
		"token":        accessToken,
		"refreshToken": refreshToken,
	})
}

// ForgotPassword Handler
func ForgotPassword(c *gin.Context) {
	var requestData struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	requestData.Email = strings.TrimSpace(requestData.Email)

	if requestData.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	token, err := services.RequestPasswordReset(requestData.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset link has been sent", "token": token})
}

// ResetPassword Handler
func ResetPassword(c *gin.Context) {
	var requestData struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if requestData.Token == "" || requestData.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token and new password are required"})
		return
	}

	err := services.ResetPassword(requestData.Token, requestData.NewPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetUserProfile Handler
func GetUserProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := repositories.GetUserByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID.Hex(),
		"name":       user.Name,
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

func Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing token"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	err = repositories.BlacklistToken(tokenString, claims.ExpiresAt.Time)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to blacklist token"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err == nil {
		_ = repositories.DeleteAllRefreshTokensForUser(objectID)
	}

	c.Status(http.StatusNoContent)
}

func RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.RefreshToken) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token required"})
		return
	}

	stored, err := repositories.FindRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	user, err := repositories.GetUserByID(stored.UserID.Hex())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	err = repositories.DeleteAllRefreshTokensForUser(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke old refresh tokens"})
		return
	}

	newAccessToken, err := utils.GenerateJWT(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	err = repositories.SaveRefreshToken(user.ID, newRefreshToken, time.Now().Add(7*24*time.Hour))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":        newAccessToken,
		"refreshToken": newRefreshToken,
	})
}
