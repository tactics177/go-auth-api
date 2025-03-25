package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/tactics177/go-auth-api/config"
	"github.com/tactics177/go-auth-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SaveRefreshToken(userID primitive.ObjectID, token string, expiresAt time.Time) error {
	collection := config.DB.Collection("refresh_tokens")

	refreshToken := models.RefreshToken{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	_, err := collection.InsertOne(context.Background(), refreshToken)
	return err
}

func FindRefreshToken(token string) (*models.RefreshToken, error) {
	collection := config.DB.Collection("refresh_tokens")

	var result models.RefreshToken
	err := collection.FindOne(context.Background(), bson.M{"token": token}).Decode(&result)
	if err != nil {
		return nil, errors.New("refresh token not found or expired")
	}

	return &result, nil
}

func DeleteAllRefreshTokensForUser(userID primitive.ObjectID) error {
	collection := config.DB.Collection("refresh_tokens")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.DeleteMany(ctx, bson.M{"user_id": userID})
	return err
}
