package transaction

import (
	"third-party/models"
	"third-party/storage"

	"go.mongodb.org/mongo-driver/bson"
)

type TransactionService interface {
	CreatePaymentTransaction(accountId, reference string, amount float64) (models.Transaction, error)
	GetPaymentTransaction(reference string) (models.Transaction, error)
}

type transactionService struct {
	Storage storage.Storage
}

func NewTransactionService(storage storage.Storage) TransactionService {
	return &transactionService{
		Storage: storage,
	}
}

func (t *transactionService) CreatePaymentTransaction(accountId, reference string, amount float64) (models.Transaction, error) {
	transaction := models.Transaction{
		AccountId: accountId,
		Reference: reference,
		Amount:    amount,
	}
	_, err := t.Storage.CreateOneRecord(&transaction)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}
func (t *transactionService) GetPaymentTransaction(reference string) (models.Transaction, error) {
	transaction := models.Transaction{}
	_, err := t.Storage.SelectOneFromDb(&transaction, bson.M{"reference": reference})
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}
