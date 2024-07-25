package wallet

import (
	"banking-app/models"
	third_party "banking-app/pkg/third-party"
	"banking-app/services/transaction"
	"banking-app/storage"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WalletService interface {
	DebitWallet(accountID primitive.ObjectID, amount float64, reference string) (float64, error)
	CreditWallet(accountID primitive.ObjectID, amount float64, reference string) (float64, error)
}

type walletService struct {
	Storage    storage.Storage
	ThirdParty third_party.ThirdPartyPkg
}

func NewWalletService(storage storage.Storage, thirdParty third_party.ThirdPartyPkg) WalletService {
	return &walletService{
		Storage:    storage,
		ThirdParty: thirdParty,
	}
}

func (w *walletService) DebitWallet(accountID primitive.ObjectID, amount float64, reference string) (float64, error) {
	transactionService := transaction.NewTransactionService(w.Storage)
	if amount <= 0 {
		return 0, fmt.Errorf("amount must be greater than zero")
	}

	if reference == "" {
		return 0, fmt.Errorf("missing transaction reference")
	}

	_, err := transactionService.GetTransaction(reference)
	if err == nil {
		return 0, fmt.Errorf("reference already exists")
	}

	wallet := &models.Wallet{}
	walletData, err := w.Storage.SelectOneFromDb(wallet, bson.M{"account_id": accountID})
	fmt.Println("----2-2-2-2-2---------", walletData, err)
	if err != nil {
		return 0, err
	}

	wallet, ok := walletData.(*models.Wallet)
	if !ok {
		return 0, fmt.Errorf("unable to cast to Wallet")
	}

	if amount > wallet.Balance {
		return 0, fmt.Errorf("insufficient funds")
	}

	wallet.Balance = wallet.Balance - amount
	walletData, err = w.Storage.UpdateRecord(wallet)
	if err != nil {
		return wallet.Balance, err
	}

	wallet, ok = walletData.(*models.Wallet)
	if !ok {
		return 0, fmt.Errorf("unable to cast to Wallet")
	}

	transaction, err := transactionService.CreateTransaction(accountID, models.DebitTransaction, amount, reference)
	if err != nil {
		return wallet.Balance, err
	}

	_, err = w.ThirdParty.CreateTransaction(accountID.String(), transaction.Reference, amount)
	if err != nil {
		return wallet.Balance, err
	}

	return wallet.Balance, nil
}

func (w *walletService) CreditWallet(accountID primitive.ObjectID, amount float64, reference string) (float64, error) {
	transactionService := transaction.NewTransactionService(w.Storage)
	if amount <= 0 {
		return 0, fmt.Errorf("amount must be greater than zero")
	}

	if reference == "" {
		return 0, fmt.Errorf("missing transaction reference")
	}

	_, err := transactionService.GetTransaction(reference)
	if err == nil {
		return 0, fmt.Errorf("reference already exists")
	}

	wallet := &models.Wallet{}
	walletData, err := w.Storage.SelectOneFromDb(wallet, bson.M{"account_id": accountID})
	if err != nil {
		return 0, err
	}

	wallet, ok := walletData.(*models.Wallet)
	if !ok {
		return 0, fmt.Errorf("unable to cast to Wallet")
	}

	wallet.Balance = wallet.Balance + amount
	walletData, err = w.Storage.UpdateRecord(wallet)
	if err != nil {
		return wallet.Balance, err
	}

	wallet, ok = walletData.(*models.Wallet)
	if !ok {
		return 0, fmt.Errorf("unable to cast to Wallet")
	}

	transaction, err := transactionService.CreateTransaction(accountID, models.CreditTransaction, amount, reference)
	if err != nil {
		return wallet.Balance, err
	}

	_, err = w.ThirdParty.CreateTransaction(accountID.String(), transaction.Reference, amount)
	if err != nil {
		return wallet.Balance, err
	}

	return wallet.Balance, nil
}
