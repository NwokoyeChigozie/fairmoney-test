package transaction

import (
	"banking-app/models"
	"banking-app/storage"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionService interface {
	CreateTransaction(accountID primitive.ObjectID, transactionType models.TransactionType, amount float64, reference string) (*models.Transaction, error)
	GetTransaction(reference string) (*models.Transaction, error)
}

type transactionService struct {
	Storage storage.Storage
}

func NewTransactionService(storage storage.Storage) TransactionService {
	return &transactionService{
		Storage: storage,
	}
}

func (t *transactionService) CreateTransaction(accountID primitive.ObjectID, transactionType models.TransactionType, amount float64, reference string) (*models.Transaction, error) {
	transaction := &models.Transaction{
		AccountId: accountID,
		Type:      transactionType,
		Reference: reference,
		Amount:    amount,
	}

	transactionData, err := t.Storage.CreateOneRecord(transaction)
	if err != nil {
		return transaction, err
	}

	transaction, ok := transactionData.(*models.Transaction)
	if !ok {
		return transaction, fmt.Errorf("unable to cast to Transaction")
	}

	return transaction, nil
}

func (t *transactionService) GetTransaction(reference string) (*models.Transaction, error) {
	transaction := &models.Transaction{}

	transactionData, err := t.Storage.SelectOneFromDb(transaction, bson.M{"reference": reference})
	if err != nil {
		return transaction, err
	}

	transaction, ok := transactionData.(*models.Transaction)
	if !ok {
		return transaction, fmt.Errorf("unable to cast to Transaction")
	}

	return transaction, nil
}
