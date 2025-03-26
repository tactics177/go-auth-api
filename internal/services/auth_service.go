package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/tactics177/go-auth-api/internal/models"
	"github.com/tactics177/go-auth-api/internal/repositories"
	"github.com/tactics177/go-auth-api/internal/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var RegisterUserFn = RegisterUser
var AuthenticateUserFn = AuthenticateUser
var RequestPasswordResetFn = RequestPasswordReset
var ResetPasswordFn = ResetPassword

// RegisterUser handles user registration
func RegisterUser(user *models.User) error {
	if !utils.ValidateEmail(user.Email) {
		return errors.New("invalid email format")
	}

	if !utils.ValidatePassword(user.Password) {
		return errors.New("password must be at least 8 characters long and contain at least one letter and one number")
	}

	existingUser, _ := repositories.GetUserByEmail(user.Email)
	if existingUser != nil {
		return errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error hashing password")
	}

	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err = repositories.CreateUser(user)
	if err != nil {
		return errors.New("error creating user")
	}

	return nil
}

// AuthenticateUser validates login credentials
func AuthenticateUser(email, password string) (string, error) {
	user, err := repositories.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := utils.GenerateJWT(*user)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func GenerateResetToken() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// RequestPasswordReset handles forgot-password requests
func RequestPasswordReset(email string) (string, error) {
	if !utils.ValidateEmail(email) {
		return "", errors.New("invalid email format")
	}

	user, err := repositories.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("email not found")
	}

	token, err := GenerateResetToken()
	if err != nil {
		return "", errors.New("failed to generate reset token")
	}

	err = repositories.SavePasswordResetToken(user.ID, token)
	if err != nil {
		return "", errors.New("failed to save reset token")
	}

	// TODO: Send reset token via email

	return token, nil
}

// ResetPassword handles resetting a user's password
func ResetPassword(token, newPassword string) error {
	if !utils.ValidatePassword(newPassword) {
		return errors.New("password must be at least 8 characters long and contain at least one letter and one number")
	}

	resetRecord, err := repositories.GetPasswordResetToken(token)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error hashing new password")
	}

	err = repositories.UpdateUserPassword(resetRecord.UserID, string(hashedPassword))
	if err != nil {
		return errors.New("failed to update password")
	}

	_ = repositories.DeletePasswordResetToken(token)

	return nil
}
