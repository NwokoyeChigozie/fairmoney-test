package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (Transaction) CollectionName() string {
	return "transactions"
}

type Transaction struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	AccountId string             `bson:"account_id" index:"account_id" json:"account_id"`
	Reference string             `bson:"reference" json:"reference"`
	Amount    float64            `bson:"amount" json:"amount"`
	CreatedAt time.Time          `bson:"created_at" json:"-"`
	UpdatedAt time.Time          `bson:"updated_at" json:"-"`
}
