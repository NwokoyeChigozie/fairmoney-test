package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (Transaction) CollectionName() string {
	return "transactions"
}

type TransactionType string

var (
	DebitTransaction  TransactionType = "debit"
	CreditTransaction TransactionType = "credit"
)

type Transaction struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	AccountId primitive.ObjectID `bson:"account_id" index:"account_id" json:"account_id"`
	Type      TransactionType    `bson:"type" json:"type"`
	Reference string             `bson:"reference" index:"account_id,unique" json:"reference"`
	Amount    float64            `bson:"amount" json:"amount"`
	CreatedAt time.Time          `bson:"created_at" json:"-"`
	UpdatedAt time.Time          `bson:"updated_at" json:"-"`
}
