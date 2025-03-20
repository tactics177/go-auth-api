package utils

import (
	"github.com/tactics177/go-auth-api/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tactics177/go-auth-api/config"
)

func GenerateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JwtSecret))
}
