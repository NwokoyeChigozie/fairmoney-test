package wallet

import (
	storage_mock "banking-app/mocks/storage"
	third_party_mock "banking-app/mocks/third_party"
	"banking-app/models"
	third_party "banking-app/pkg/third-party"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestWalletService_CreditWallet(t *testing.T) {
	// id := primitive.NewObjectID()
	accountid := primitive.NewObjectID()
	uid, _ := uuid.NewV4()
	reference := uid.String()
	amount := 100.00
	// currentTime := time.Now()
	tType := models.CreditTransaction
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storage_mock.NewMockStorage(ctrl)
	thirdParty := third_party_mock.NewMockThirdPartyPkg(ctrl)
	walletService := NewWalletService(mockStorage, thirdParty)

	transaction := models.Transaction{
		AccountId: accountid,
		Type:      tType,
		Reference: reference,
		Amount:    amount,
	}
	thirdPartyModel := &third_party.ThirdPartyPackageResponse{
		AccountId: accountid.String(),
		Reference: reference,
		Amount:    amount,
	}

	t.Run("successful credit transaction", func(t *testing.T) {
		existingtransaction := models.Transaction{}
		mockStorage.EXPECT().SelectOneFromDb(&existingtransaction, bson.M{"reference": reference}).Return(&existingtransaction, mongo.ErrNoDocuments)
		mockStorage.EXPECT().SelectOneFromDb(&models.Wallet{}, bson.M{"account_id": accountid}).Return(&models.Wallet{}, nil)
		mockStorage.EXPECT().UpdateRecord(&models.Wallet{
			Balance: amount,
		}).Return(&models.Wallet{Balance: amount}, nil)
		mockStorage.EXPECT().CreateOneRecord(&transaction).Return(&models.Transaction{Reference: reference}, nil)
		thirdParty.EXPECT().CreateTransaction(accountid.String(), reference, amount).Return(thirdPartyModel, nil)

		balance, err := walletService.CreditWallet(accountid, amount, reference)
		assert.Equal(t, err, nil)
		assert.Equal(t, balance, amount)
	})
	t.Run("amount less than zero", func(t *testing.T) {
		creditAmount := -100.00

		balance, err := walletService.CreditWallet(accountid, creditAmount, reference)
		assert.NotEqual(t, err, nil)
		assert.Equal(t, balance, 0.00)
		assert.Contains(t, err.Error(), "amount must be greater than zero")
	})

	t.Run("invalid account id", func(t *testing.T) {
		invalidAccountId := primitive.NewObjectID()
		existingtransaction := models.Transaction{}
		mockStorage.EXPECT().SelectOneFromDb(&existingtransaction, bson.M{"reference": reference}).Return(&existingtransaction, mongo.ErrNoDocuments)
		mockStorage.EXPECT().SelectOneFromDb(&models.Wallet{}, bson.M{"account_id": invalidAccountId}).Return(&models.Wallet{}, mongo.ErrNoDocuments)

		balance, err := walletService.CreditWallet(invalidAccountId, amount, reference)
		assert.NotEqual(t, err, nil)
		assert.Equal(t, err, mongo.ErrNoDocuments)
		assert.Equal(t, balance, 0.00)
	})

	t.Run("missing transaction reference", func(t *testing.T) {

		balance, err := walletService.CreditWallet(accountid, amount, "")
		assert.NotEqual(t, err, nil)
		assert.Contains(t, err.Error(), "missing transaction reference")
		assert.Equal(t, balance, 0.00)
	})

	t.Run("reference already exists", func(t *testing.T) {
		existingtransaction := models.Transaction{}
		mockStorage.EXPECT().SelectOneFromDb(&existingtransaction, bson.M{"reference": reference}).Return(&existingtransaction, nil)

		balance, err := walletService.CreditWallet(accountid, amount, reference)
		assert.NotEqual(t, err, nil)
		assert.Equal(t, balance, 0.00)
		assert.Contains(t, err.Error(), "reference already exists")
	})

}

func TestWalletService_DebitWallet(t *testing.T) {
	// id := primitive.NewObjectID()
	accountid := primitive.NewObjectID()
	uid, _ := uuid.NewV4()
	reference := uid.String()
	amount := 50.00
	// currentTime := time.Now()
	ctrl := gomock.NewController(t)
	tType := models.DebitTransaction
	defer ctrl.Finish()

	mockStorage := storage_mock.NewMockStorage(ctrl)
	thirdParty := third_party_mock.NewMockThirdPartyPkg(ctrl)
	walletService := NewWalletService(mockStorage, thirdParty)

	transaction := models.Transaction{
		AccountId: accountid,
		Type:      tType,
		Reference: reference,
		Amount:    amount,
	}
	thirdPartyModel := &third_party.ThirdPartyPackageResponse{
		AccountId: accountid.String(),
		Reference: reference,
		Amount:    amount,
	}

	t.Run("successful debit transaction", func(t *testing.T) {
		baseWallet := models.Wallet{}
		currentBalance := 100.00
		existingtransaction := models.Transaction{}
		mockStorage.EXPECT().SelectOneFromDb(&existingtransaction, bson.M{"reference": reference}).Return(&existingtransaction, mongo.ErrNoDocuments)
		mockStorage.EXPECT().SelectOneFromDb(&baseWallet, bson.M{"account_id": accountid}).Return(&models.Wallet{Balance: currentBalance}, nil)
		mockStorage.EXPECT().UpdateRecord(&models.Wallet{
			Balance: amount,
		}).Return(&models.Wallet{Balance: currentBalance - amount}, nil)
		mockStorage.EXPECT().CreateOneRecord(&transaction).Return(&transaction, nil)
		thirdParty.EXPECT().CreateTransaction(accountid.String(), reference, amount).Return(thirdPartyModel, nil)

		balance, err := walletService.DebitWallet(accountid, amount, reference)
		assert.Equal(t, err, nil)
		assert.Equal(t, balance, amount)
	})

	t.Run("amount less than zero", func(t *testing.T) {
		debitAmount := -100.00

		balance, err := walletService.DebitWallet(accountid, debitAmount, reference)
		assert.NotEqual(t, err, nil)
		assert.Equal(t, balance, 0.00)
		assert.Contains(t, err.Error(), "amount must be greater than zero")
	})

	t.Run("missing transaction reference", func(t *testing.T) {
		balance, err := walletService.DebitWallet(accountid, amount, "")
		assert.NotEqual(t, err, nil)
		assert.Contains(t, err.Error(), "missing transaction reference")
		assert.Equal(t, balance, 0.00)
	})

	t.Run("invalid account id", func(t *testing.T) {
		invalidAccountId := primitive.NewObjectID()
		existingtransaction := models.Transaction{}
		mockStorage.EXPECT().SelectOneFromDb(&existingtransaction, bson.M{"reference": reference}).Return(&existingtransaction, mongo.ErrNoDocuments)
		mockStorage.EXPECT().SelectOneFromDb(&models.Wallet{}, bson.M{"account_id": invalidAccountId}).Return(&models.Wallet{}, mongo.ErrNoDocuments)

		balance, err := walletService.DebitWallet(invalidAccountId, amount, reference)
		assert.NotEqual(t, err, nil)
		assert.Equal(t, err, mongo.ErrNoDocuments)
		assert.Equal(t, balance, 0.00)
	})

	t.Run("insufficient balance", func(t *testing.T) {
		baseWallet := models.Wallet{}
		currentBalance := 100.00
		debitAmount := 200.00
		existingtransaction := models.Transaction{}
		mockStorage.EXPECT().SelectOneFromDb(&existingtransaction, bson.M{"reference": reference}).Return(&existingtransaction, mongo.ErrNoDocuments)
		mockStorage.EXPECT().SelectOneFromDb(&baseWallet, bson.M{"account_id": accountid}).Return(&models.Wallet{Balance: currentBalance}, nil)

		balance, err := walletService.DebitWallet(accountid, debitAmount, reference)
		assert.NotEqual(t, err, nil)
		assert.Equal(t, balance, 0.00)
		assert.Contains(t, err.Error(), "insufficient funds")
	})

	t.Run("reference already exists", func(t *testing.T) {
		existingtransaction := models.Transaction{}
		mockStorage.EXPECT().SelectOneFromDb(&existingtransaction, bson.M{"reference": reference}).Return(&existingtransaction, nil)

		balance, err := walletService.DebitWallet(accountid, amount, reference)
		assert.NotEqual(t, err, nil)
		assert.Equal(t, balance, 0.00)
		assert.Contains(t, err.Error(), "reference already exists")
	})

}
