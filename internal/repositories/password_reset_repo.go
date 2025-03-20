package repositories

import (
	"context"
	"errors"
	"github.com/tactics177/go-auth-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
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

func GetPasswordResetToken(token string) (*models.PasswordReset, error) {
	resetCollection := config.DB.Collection("password_resets")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var resetRecord models.PasswordReset
	err := resetCollection.FindOne(ctx, bson.M{"token": token}).Decode(&resetRecord)
	if err != nil {
		return nil, errors.New("invalid or expired reset token")
	}

	if time.Now().After(resetRecord.ExpiresAt) {
		_ = DeletePasswordResetToken(token)
		return nil, errors.New("reset token has expired")
	}

	return &resetRecord, nil
}

func DeletePasswordResetToken(token string) error {
	resetCollection := config.DB.Collection("password_resets")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := resetCollection.DeleteOne(ctx, bson.M{"token": token})
	return err
}
