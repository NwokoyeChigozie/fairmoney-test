package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (User) CollectionName() string {
	return "users"
}

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Email     string             `bson:"email" json:"email"`
	CreatedAt time.Time          `bson:"created_at" json:"-"`
	UpdatedAt time.Time          `bson:"updated_at" json:"-"`
}
