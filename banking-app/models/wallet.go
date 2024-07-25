package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (Wallet) CollectionName() string {
	return "wallets"
}

type Wallet struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	AccountId primitive.ObjectID `bson:"account_id" json:"account_id"`
	Balance   float64            `bson:"amount" json:"amount"`
	CreatedAt time.Time          `bson:"created_at" json:"-"`
	UpdatedAt time.Time          `bson:"updated_at" json:"-"`
}
