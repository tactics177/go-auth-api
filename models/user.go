package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` //to be generated automatically by mongodb
	Name      string             `bson:"name"`          //always required
	Email     string             `bson:"email"`         //always required
	Password  string             `bson:"password"`      //always required
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
}
