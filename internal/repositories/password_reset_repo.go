package repositories

import (
	"context"
	"github.com/tactics177/go-auth-api/internal/models"
	"time"

	"github.com/tactics177/go-auth-api/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SavePasswordResetToken(userID primitive.ObjectID, token string) error {
	resetCollection := config.DB.Collection("password_resets")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resetRecord := models.PasswordReset{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Minute * 15),
	}

	_, err := resetCollection.InsertOne(ctx, resetRecord)
	return err
}
