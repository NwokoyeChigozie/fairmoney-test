package http_handler

import (
	storage_mock "banking-app/mocks/storage"
	third_party_mock "banking-app/mocks/third_party"
	"banking-app/models"
	third_party "banking-app/pkg/third-party"
	"banking-app/services/wallet"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestWalletHandler_CreditWallet(t *testing.T) {
	accountid := primitive.NewObjectID()
	uid, _ := uuid.NewV4()
	reference := uid.String()
	amount := 100.00
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storage_mock.NewMockStorage(ctrl)
	thirdParty := third_party_mock.NewMockThirdPartyPkg(ctrl)
	mockWalletService := wallet.NewWalletService(mockStorage, thirdParty)
	walletHandler := NewWalletHandler(mockWalletService)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	route := "/bank-api/credit-wallet"
	r.POST(route, walletHandler.CreditWallet)
	transaction := models.Transaction{
		AccountId: accountid,
		Type:      models.CreditTransaction,
		Reference: reference,
		Amount:    amount,
	}

	thirdPartyModel := &third_party.ThirdPartyPackageResponse{
		AccountId: accountid.String(),
		Reference: reference,
		Amount:    amount,
	}
	t.Run("successful wallet credit", func(t *testing.T) {

		existingtransaction := models.Transaction{}
		mockStorage.EXPECT().SelectOneFromDb(&existingtransaction, bson.M{"reference": reference}).Return(&existingtransaction, mongo.ErrNoDocuments)
		mockStorage.EXPECT().SelectOneFromDb(&models.Wallet{}, bson.M{"account_id": accountid}).Return(&models.Wallet{}, nil)
		mockStorage.EXPECT().UpdateRecord(&models.Wallet{
			Balance: amount,
		}).Return(&models.Wallet{Balance: amount}, nil)
		mockStorage.EXPECT().CreateOneRecord(&transaction).Return(&models.Transaction{Reference: reference}, nil)
		thirdParty.EXPECT().CreateTransaction(accountid.String(), reference, amount).Return(thirdPartyModel, nil)

		bodyModel := map[string]interface{}{
			"amount":     amount,
			"account_id": accountid,
			"reference":  reference,
		}

		body, err := json.Marshal(bodyModel)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}
		req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"balance":100`)
	})

	t.Run("account id is required", func(t *testing.T) {
		bodyModel := map[string]interface{}{
			"amount":    amount,
			"reference": reference,
		}

		body, err := json.Marshal(bodyModel)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}
		req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `'AccountId' failed on the 'required' tag"`)
	})

	t.Run("reference is required", func(t *testing.T) {
		bodyModel := map[string]interface{}{
			"amount":     amount,
			"account_id": accountid,
		}

		body, err := json.Marshal(bodyModel)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}
		req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		fmt.Println("oiuqebvubouqbeobqeoqe", w.Body.String())
		assert.Contains(t, w.Body.String(), `'Reference' failed on the 'required' tag"`)
	})

	t.Run("amount is required", func(t *testing.T) {
		bodyModel := map[string]interface{}{
			"account_id": accountid,
			"reference":  reference,
		}

		body, err := json.Marshal(bodyModel)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}
		req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		fmt.Println("oiuqebvubouqbeobqeoqe", w.Body.String())
		assert.Contains(t, w.Body.String(), `'Amount' failed on the 'required' tag"`)
	})

	t.Run("amount must be greater than 0", func(t *testing.T) {
		bodyModel := map[string]interface{}{
			"amount":     -100,
			"account_id": accountid,
			"reference":  reference,
		}

		body, err := json.Marshal(bodyModel)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}
		req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		fmt.Println("oiuqebvubouqbeobqeoqe", w.Body.String())
		assert.Contains(t, w.Body.String(), `amount must be greater than zero`)
	})

	t.Run("reference already exists", func(t *testing.T) {

		existingtransaction := models.Transaction{}
		mockStorage.EXPECT().SelectOneFromDb(&existingtransaction, bson.M{"reference": reference}).Return(&existingtransaction, nil)

		bodyModel := map[string]interface{}{
			"amount":     amount,
			"account_id": accountid,
			"reference":  reference,
		}

		body, err := json.Marshal(bodyModel)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}
		req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `reference already exists`)
	})

	t.Run("invalid account id", func(t *testing.T) {

		existingtransaction := models.Transaction{}
		mockStorage.EXPECT().SelectOneFromDb(&existingtransaction, bson.M{"reference": reference}).Return(&existingtransaction, mongo.ErrNoDocuments)
		mockStorage.EXPECT().SelectOneFromDb(&models.Wallet{}, bson.M{"account_id": accountid}).Return(&models.Wallet{}, mongo.ErrNoDocuments)

		bodyModel := map[string]interface{}{
			"amount":     amount,
			"account_id": accountid,
			"reference":  reference,
		}

		body, err := json.Marshal(bodyModel)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}
		req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), mongo.ErrNoDocuments.Error())
	})

}

func TestWalletHandler_DebitWallet(t *testing.T) {
	accountid := primitive.NewObjectID()
	uid, _ := uuid.NewV4()
	reference := uid.String()
	amount := 100.00
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storage_mock.NewMockStorage(ctrl)
	thirdParty := third_party_mock.NewMockThirdPartyPkg(ctrl)
	mockWalletService := wallet.NewWalletService(mockStorage, thirdParty)
	walletHandler := NewWalletHandler(mockWalletService)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	route := "/bank-api/debit-wallet"
	r.POST(route, walletHandler.DebitWallet)
	transaction := models.Transaction{
		AccountId: accountid,
		Type:      models.DebitTransaction,
		Reference: reference,
		Amount:    amount,
	}

	thirdPartyModel := &third_party.ThirdPartyPackageResponse{
		AccountId: accountid.String(),
		Reference: reference,
		Amount:    amount,
	}
	t.Run("successful wallet debit", func(t *testing.T) {

		baseWallet := models.Wallet{}
		currentBalance := 200.00
		existingtransaction := models.Transaction{}
		mockStorage.EXPECT().SelectOneFromDb(&existingtransaction, bson.M{"reference": reference}).Return(&existingtransaction, mongo.ErrNoDocuments)
		mockStorage.EXPECT().SelectOneFromDb(&baseWallet, bson.M{"account_id": accountid}).Return(&models.Wallet{Balance: currentBalance}, nil)
		mockStorage.EXPECT().UpdateRecord(&models.Wallet{
			Balance: amount,
		}).Return(&models.Wallet{Balance: currentBalance - amount}, nil)
		mockStorage.EXPECT().CreateOneRecord(&transaction).Return(&transaction, nil)
		thirdParty.EXPECT().CreateTransaction(accountid.String(), reference, amount).Return(thirdPartyModel, nil)

		bodyModel := map[string]interface{}{
			"amount":     amount,
			"account_id": accountid,
			"reference":  reference,
		}

		body, err := json.Marshal(bodyModel)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}
		req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"balance":100`)
	})
	t.Run("insufficient balance", func(t *testing.T) {

		baseWallet := models.Wallet{}
		currentBalance := 50.00
		existingtransaction := models.Transaction{}
		mockStorage.EXPECT().SelectOneFromDb(&existingtransaction, bson.M{"reference": reference}).Return(&existingtransaction, mongo.ErrNoDocuments)
		mockStorage.EXPECT().SelectOneFromDb(&baseWallet, bson.M{"account_id": accountid}).Return(&models.Wallet{Balance: currentBalance}, nil)

		bodyModel := map[string]interface{}{
			"amount":     amount,
			"account_id": accountid,
			"reference":  reference,
		}

		body, err := json.Marshal(bodyModel)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}
		req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `insufficient funds`)
	})

	t.Run("account id is required", func(t *testing.T) {
		bodyModel := map[string]interface{}{
			"amount":    amount,
			"reference": reference,
		}

		body, err := json.Marshal(bodyModel)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}
		req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `'AccountId' failed on the 'required' tag"`)
	})

	t.Run("reference is required", func(t *testing.T) {
		bodyModel := map[string]interface{}{
			"amount":     amount,
			"account_id": accountid,
		}

		body, err := json.Marshal(bodyModel)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}
		req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		fmt.Println("oiuqebvubouqbeobqeoqe", w.Body.String())
		assert.Contains(t, w.Body.String(), `'Reference' failed on the 'required' tag"`)
	})

	t.Run("amount is required", func(t *testing.T) {
		bodyModel := map[string]interface{}{
			"account_id": accountid,
			"reference":  reference,
		}

		body, err := json.Marshal(bodyModel)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}
		req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		fmt.Println("oiuqebvubouqbeobqeoqe", w.Body.String())
		assert.Contains(t, w.Body.String(), `'Amount' failed on the 'required' tag"`)
	})

	t.Run("amount must be greater than 0", func(t *testing.T) {
		bodyModel := map[string]interface{}{
			"amount":     -100,
			"account_id": accountid,
			"reference":  reference,
		}

		body, err := json.Marshal(bodyModel)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}
		req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		fmt.Println("oiuqebvubouqbeobqeoqe", w.Body.String())
		assert.Contains(t, w.Body.String(), `amount must be greater than zero`)
	})

	t.Run("reference already exists", func(t *testing.T) {

		existingtransaction := models.Transaction{}
		mockStorage.EXPECT().SelectOneFromDb(&existingtransaction, bson.M{"reference": reference}).Return(&existingtransaction, nil)

		bodyModel := map[string]interface{}{
			"amount":     amount,
			"account_id": accountid,
			"reference":  reference,
		}

		body, err := json.Marshal(bodyModel)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}
		req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `reference already exists`)
	})

	t.Run("invalid account id", func(t *testing.T) {

		existingtransaction := models.Transaction{}
		mockStorage.EXPECT().SelectOneFromDb(&existingtransaction, bson.M{"reference": reference}).Return(&existingtransaction, mongo.ErrNoDocuments)
		mockStorage.EXPECT().SelectOneFromDb(&models.Wallet{}, bson.M{"account_id": accountid}).Return(&models.Wallet{}, mongo.ErrNoDocuments)

		bodyModel := map[string]interface{}{
			"amount":     amount,
			"account_id": accountid,
			"reference":  reference,
		}

		body, err := json.Marshal(bodyModel)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}
		req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), mongo.ErrNoDocuments.Error())
	})

}
