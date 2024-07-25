package transaction

import (
	"testing"
	storage_mock "third-party/mocks/storage"
	"third-party/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestWalletHandler_Credit(t *testing.T) {
	accountid := "7727UBEUBBEI9"
	reference := "YYUS27UBEUBBEI9"
	amount := 500.00
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorage := storage_mock.NewMockStorage(ctrl)
	service := NewTransactionService(mockStorage)
	newTransaction := models.Transaction{
		AccountId: accountid,
		Reference: reference,
		Amount:    amount,
	}
	// transaction := models.Transaction{
	// 	ID:        primitive.NewObjectID(),
	// 	AccountId: accountid,
	// 	Reference: reference,
	// 	Amount:    amount,
	// 	CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now(),
	// }

	t.Run("successfully create transaction", func(t *testing.T) {
		mockStorage.EXPECT().CreateOneRecord(&newTransaction).Return(newTransaction, nil)

		_, err := service.CreatePaymentTransaction(accountid, reference, amount)
		assert.Equal(t, err, nil)

	})

	t.Run("successfully get transaction", func(t *testing.T) {
		eTransaction := models.Transaction{}
		mockStorage.EXPECT().SelectOneFromDb(&eTransaction, bson.M{"reference": reference}).Return(newTransaction, nil)

		_, err := service.GetPaymentTransaction(reference)
		assert.Equal(t, err, nil)

	})

}
