package transaction

import (
	storage_mock "banking-app/mocks/storage"
	"banking-app/models"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTransactionService_CreateTransaction(t *testing.T) {
	// id := primitive.NewObjectID()
	accountid := primitive.NewObjectID()
	uid, _ := uuid.NewV4()
	reference := uid.String()
	amount := 500.00
	tType := models.CreditTransaction
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storage_mock.NewMockStorage(ctrl)
	service := NewTransactionService(mockStorage)

	newTransaction := models.Transaction{
		AccountId: accountid,
		Type:      tType,
		Reference: reference,
		Amount:    amount,
	}

	t.Run("successfully credit transaction", func(t *testing.T) {
		mockStorage.EXPECT().CreateOneRecord(&newTransaction).Return(&newTransaction, nil)

		_, err := service.CreateTransaction(accountid, models.CreditTransaction, amount, reference)
		assert.Equal(t, err, nil)
	})

}
