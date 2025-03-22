package handlers

import (
	"github.com/tactics177/go-auth-api/internal/models"
	"github.com/tactics177/go-auth-api/internal/repositories"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tactics177/go-auth-api/internal/services"
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

	token, err := services.AuthenticateUser(strings.TrimSpace(loginData.Email), strings.TrimSpace(loginData.Password))
	if err != nil {
		switch err.Error() {
		case "invalid email format":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "failed to generate token":
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
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

// Logout Handler
func Logout(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
