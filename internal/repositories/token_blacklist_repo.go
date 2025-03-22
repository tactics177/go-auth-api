package repositories

import (
	"context"
	"time"

	"github.com/tactics177/go-auth-api/config"
	"github.com/tactics177/go-auth-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

func BlacklistToken(token string, expiresAt time.Time) error {
	collection := config.DB.Collection("blacklisted_tokens")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, models.BlacklistedToken{
		Token:     token,
		ExpiresAt: expiresAt,
	})
	return err
}

func IsTokenBlacklisted(token string) (bool, error) {
	collection := config.DB.Collection("blacklisted_tokens")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, bson.M{"token": token})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
