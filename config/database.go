package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var JwtSecret string

func ConnectDB() {
	err := godotenv.Load()
	// .env for local development, environment variables for cloud deployment
	if err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not found in environment")
	}

	JwtSecret = os.Getenv("JWT_SECRET")
	if JwtSecret == "" {
		log.Fatal("JWT_SECRET not found in environment")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB ping failed:", err)
	}

	fmt.Println("Connected to MongoDB Atlas successfully!")

	DB = client.Database("go-auth-api")

	createIndexes()
	fmt.Println("JWT Secret loaded successfully!")
}

func createIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := DB.Collection("users")
	resetCollection := DB.Collection("password_resets")

	emailIndex := mongo.IndexModel{
		Keys:    bson.D{{"email", 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := userCollection.Indexes().CreateOne(ctx, emailIndex)
	if err != nil {
		log.Fatal("Error creating email index:", err)
	}

	// TTL index to auto-delete expired password reset tokens
	resetIndex := mongo.IndexModel{
		Keys:    bson.D{{"expires_at", 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	}
	_, err = resetCollection.Indexes().CreateOne(ctx, resetIndex)
	if err != nil {
		log.Fatal("Error creating password reset TTL index:", err)
	}

	// TTL index on blacklisted_tokens
	blacklistCollection := DB.Collection("blacklisted_tokens")

	ttlIndexBlacklistedToken := mongo.IndexModel{
		Keys:    bson.D{{Key: "expires_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	}
	_, err = blacklistCollection.Indexes().CreateOne(ctx, ttlIndexBlacklistedToken)
	if err != nil {
		log.Fatal("Error creating TTL index on blacklisted_tokens:", err)
	}

	// TTL index on refresh_tokens
	refreshTokens := DB.Collection("refresh_tokens")

	ttlIndexRefreshToken := mongo.IndexModel{
		Keys: bson.D{{Key: "expires_at", Value: 1}},
		Options: options.Index().
			SetExpireAfterSeconds(0),
	}

	_, err = refreshTokens.Indexes().CreateOne(ctx, ttlIndexRefreshToken)
	if err != nil {
		log.Fatal("Error creating TTL index for refresh_tokens:", err)
	}

	fmt.Println("Indexes created successfully!")
}
